package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/runeimp/divvy"
)

const (
	appName    = "Divvy Up the LOOT"
	appVersion = "1.0.0"
	cliName    = "divvy"
	usage      = `Usage: %s [OPTIONS] CSV_FILE

OPTIONS:
`
	usageProlog = `
CSV_FILE expects a CSV with the first column being the players and the rest
of the columns being token counts per player with the header being the name
of the prize.

Example CSV:
------------
Players,"Gold Watch","Silver Necklace",Dagger
"Lumpy Thumpkin",4,,2
"Scarlett Jewels",,4,3
"Garret Theivington",,,9
`
)

var (
	appLabel      = fmt.Sprintf("%s v%s", appName, appVersion)
	appStart      string
	args          []string
	binary        string
	binPath       string
	csvInputName  string
	csvOutputName = "DivvyUpTheLoot"
	sigs          = make(chan os.Signal, 1)
)

func main() {
	start := time.Now()
	appStart = start.Format("2006-01-02_150405_MST")

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		msg := fmt.Sprintf("SIGNAL: %q", sig.String())
		fmt.Println(msg)
	}()

	csvPtr := flag.Bool("csv", false, "Create a CSV file with the output")
	helpPtr := flag.Bool("help", false, "Display this help info")
	lootDataPtr := flag.Bool("loot", false, "Display loot data")
	lootLabelPtr := flag.String("loot-label", "Loot", "Set the loot label")
	noWinnerDataPtr := flag.Bool("no-winners", false, "Don't display winner data")
	pickPtr := flag.Int("pick", 0, "Max picks from list. Use -1 to randomize the list.")
	playerLabelPtr := flag.String("player-label", "", "Set the player label")
	tokenLabelPtr := flag.String("token-label", "Tokens", "Set the loot label")
	versionPtr := flag.Bool("version", false, "Display version info")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usage, cliName)

		flag.VisitAll(func(f *flag.Flag) {
			optionName := fmt.Sprintf("-%s", f.Name)
			if len(f.DefValue) == 0 {
				fmt.Fprintf(flag.CommandLine.Output(), "  %-13s  %s (default: 1st field of the csv header)\n", optionName, f.Usage)
			} else {
				fmt.Fprintf(flag.CommandLine.Output(), "  %-13s  %s (default: %v)\n", optionName, f.Usage, f.DefValue)
			}
		})
		fmt.Fprintf(flag.CommandLine.Output(), usageProlog)
	}

	flag.Parse()

	if *helpPtr {
		flag.Usage()
		os.Exit(0)
	}
	if *versionPtr {
		fmt.Println(appLabel)
		os.Exit(0)
	}

	args = flag.Args()

	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println(appLabel)

	if len(args) > 0 {
		csvInputName = args[0]
	}

	binary = filepath.Base(os.Args[0])
	binPath = filepath.Dir(os.Args[0])
	// fmt.Printf("binary: %q\n", binary)
	// fmt.Printf("binPath: %q\n", binPath)
	csvOutputName = fmt.Sprintf("%s_%s.csv", csvOutputName, appStart)
	// csvOutputName = fmt.Sprintf("%s/%s_%s.csv", binPath, csvOutputName, appStart)

	err := divvy.ReadCSV(csvInputName, *pickPtr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
	if *noWinnerDataPtr == false {
		winners := divvy.GetWinners()
		if *pickPtr == 0 {
			for _, winner := range winners {
				fmt.Printf("%q got %q\n", winner.Name, winner.Loot)
			}
		} else {
			for _, winner := range winners {
				fmt.Printf("%q\n", winner.Name)
			}
		}
	}

	if *lootDataPtr {
		for _, loot := range divvy.GetPrizeList() {
			if loot.Tokens == 1 {
				fmt.Printf("%s %q had 1 token\n", *lootLabelPtr, loot.Name)
			} else {
				fmt.Printf("%s %q had %d tokens\n", *lootLabelPtr, loot.Name, loot.Tokens)
			}
		}
	}

	if *csvPtr {
		err = divvy.WriteCSV(csvOutputName, *playerLabelPtr, *lootLabelPtr, *tokenLabelPtr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}
	}
}
