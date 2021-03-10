package divvy

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Player is the base data type to divvy up
type Player struct {
	Loot string
	Name string
}

// Prize is the base data type for the PrizeList
type Prize struct {
	Name   string
	Tokens int
}

var (
	csvHeaders []string
	keys       []string
	loot       []Prize
	prizes     map[string]int
	rowCount   uint
	theList    []string
	winners    []Player
)

//
// Local Functions
//
func init() {
	prizes = make(map[string]int)
}

func sortNumericStrings(i, j int) bool {
	iStr := keys[i]
	jStr := keys[j]

	re := regexp.MustCompile(`(\D*)(\d*)(.*)`)
	iReg := re.FindStringSubmatch(iStr)
	jReg := re.FindStringSubmatch(jStr)

	iPrefix := iReg[1]
	iDigit := 0
	if len(iReg[2]) > 0 {
		iDigit, _ = strconv.Atoi(iReg[2])
	}
	iSuffix := iReg[1]

	jPrefix := jReg[1]
	jDigit := 0
	if len(jReg[2]) > 0 {
		jDigit, _ = strconv.Atoi(jReg[2])
	}
	jSuffix := jReg[1]

	if iPrefix < jPrefix {
		return true
	}

	if iDigit < jDigit {
		return true
	}

	return iSuffix < jSuffix
}

func updateTheList(w, l string) {
	newList := []string{}
	for _, item := range theList {
		listParts := strings.Split(item, "::")
		winner := listParts[0]
		listPrize := listParts[1]
		if winner != w && listPrize != l {
			newList = append(newList, item)
		}
	}
	theList = []string{}
	for _, item := range newList {
		theList = append(theList, item)
	}
}

//
// Methods
//

// GetPrizeList prints the prizes and their token counts
func GetPrizeList() []Prize {
	keys = make([]string, 0, len(prizes))
	for k := range prizes {
		keys = append(keys, k)
	}
	sort.Slice(keys, sortNumericStrings)

	for _, k := range keys {
		loot = append(loot, Prize{Name: k, Tokens: prizes[k]})
	}

	return loot
}

// GetWinners returns a list of winners and their loot
func GetWinners() []Player {
	lootCount := len(csvHeaders) - 1

	if len(winners) > 0 {
		return winners
	}

	// results = append(results, strconv.Itoa(lootCount))
	for i := 0; i < lootCount; i++ {
		max := len(theList) - 1
		rand.Seed(time.Now().UnixNano())
		randy := rand.Intn(max)
		// log.Printf("max: %d | randy: %d\n", max, randy)
		item := theList[randy]
		parts := strings.Split(item, "::")
		n := parts[0]
		l := parts[1]
		player := Player{Loot: l, Name: n}
		// log.Printf("name: %s | loot: %q\n", n, l)
		winners = append(winners, player)

		updateTheList(n, l)
	}

	return winners
}

// ReadCSV parses a CSV for loot data
func ReadCSV(filepath string) error {
	log.Printf("divvy.ReadCSV() | filepath: %q\n", filepath)

	if len(filepath) == 0 {
		return errors.New("filepath name is empty")
	}

	fmt.Printf("Reading %q\nProgress: ", filepath)

	var (
		endRow uint
	)

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Printf("Error trying to read: %q", filepath)
		os.Exit(2)
	}

	header := true

	csvRecords := csv.NewReader(bytes.NewReader(data))
	for {
		rowCount++
	SkipRow:
		record, err := csvRecords.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("CSV Read Error: %s", err.Error())
		}

		if header {
			header = false
			rowCount--
			for _, field := range record {
				csvHeaders = append(csvHeaders, strings.TrimSpace(field))
			}
			// log.Printf("divvy.readCSV() | csvHeaders = %q\n", csvHeaders)
			continue
		}

		player := ""
		for i, field := range record {
			field = strings.TrimSpace(field)
			if i == 0 {
				if field == "Totals" {
					goto SkipRow
				}
				player = field
				continue
			}
			if len(field) > 0 {
				loot := csvHeaders[i]
				playerPrize := player + "::" + loot
				tokens, err := strconv.Atoi(field)
				if err != nil {
					log.Fatalln(err)
				}
				for j := 0; j < tokens; j++ {
					prizes[loot]++
					theList = append(theList, playerPrize)
				}
			}
		}

		if err != nil {
			fmt.Print("!")
		} else {
			fmt.Print(".")
		}

		if endRow > 0 && rowCount >= endRow {
			break
		}
	}

	fmt.Println()

	return nil
}

// WriteCSV creates the divvy results to a CSV file
func WriteCSV(filepath, playerLabel, lootLabel, tokenLabel string) error {
	log.Printf("divvy.WriteCSV() | filepath: %q\n", filepath)

	GetWinners()
	GetPrizeList()

	winnerCount := len(winners)

	csvData := [][]string{{playerLabel, lootLabel, lootLabel, tokenLabel}}

	record := []string{}

	for i, l := range loot {
		tokens := strconv.Itoa(l.Tokens)
		if i < winnerCount {
			winner := winners[i]
			record = []string{winner.Name, winner.Loot, l.Name, tokens}
		} else {
			record = []string{"", "", l.Name, tokens}
		}

		csvData = append(csvData, record)
	}

	csvFile, err := os.Create(filepath)
	if err != nil {
		log.Printf("CSV Creation Error: %s", err.Error())
		return err
	}

	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	for _, row := range csvData {
		err := writer.Write(row)
		if err != nil {
			log.Printf("CSV Row Writing Error: %s", err.Error())
			return err
		}
	}

	return nil
}
