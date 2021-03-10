Divvy Up the LOOT v1.0.0
========================

Command line tool for facilitating giveaways, lotteries, and raffles via CSV import and export.


Features
--------

* Defaults to showing progress and winners
* Can dissable display of winners
* Can display loot list and associated tokens as well with `-loot` or `--loot`
* Can output winners and loot list to CSV with `-csv` or `--csv`
* Can pick a fixed number of winners from a single column CSV with `-pick #` or `--pick #` where `#` is number of 

Rational
--------

My wife asked me about raffle software as one of her Discord servers was having issues finding a good solution for a giveaway they where planning. They'd already created a Google Sheet with the data they wanted to use for the giveaway so I built this (with a few additions) to take care of their situation. I also looked into command line options and couldn't find any so I created Divvy Up the LOOT. :smiling_imp:

I'm also hoping this becomes a useful command line tool for others as well. :innocent:


Examples
--------

### Get Divvy Version

```
$ divvy -version
Divvy Up the LOOT v1.0.0
```

### Get Divvy Help

```
$ divvy -h
Usage: divvy [OPTIONS] CSV_FILE

OPTIONS:
  -csv           Create a CSV file with the output (default: false)
  -help          Display this help info (default: false)
  -loot          Display loot data (default: false)
  -loot-label    Set the loot label (default: Loot)
  -no-winners    Don't display winner data (default: false)
  -pick          Max picks from list (default: 0)
  -player-label  Set the player label (default: 1st field of the csv header)
  -token-label   Set the loot label (default: Tokens)
  -version       Display version info (default: false)

CSV_FILE expects a CSV with the first column being the players and the rest
of the columns being token counts per player with the header being the name
of the prize.

Example CSV:
------------
Players,"Gold Watch","Silver Necklace",Dagger
"Lumpy Thumpkin",4,,2
"Scarlett Jewels",,4,3
"Garret Theivington",,,9
```

### Single Gift Giveaway

#### CSV Input

##### `cat random.csv`

```csv
Partisipants
"Captain Hook"
"Pegleg Peet"
"The Pirate Captain"
"The Surprisingly Curvaceous Pirate"
```

#### Run Divvy

```
$ divvy -csv -pick 1 random.csv
Divvy Up the LOOT v1.0.0
Reading "random.csv"
Progress: ....
"The Pirate Captain"
```

### Randomize List

#### CSV Input

##### `cat random.csv`

```csv
Partisipants
"Captain Hook"
"Pegleg Peet"
"The Pirate Captain"
"The Surprisingly Curvaceous Pirate"
```

#### Run Divvy

```
$ -pick -1 random.csv
Divvy Up the LOOT v1.0.0
Reading "random.csv"
Progress: ....
"The Surprisingly Curvaceous Pirate"
"Captain Hook"
"The Pirate Captain"
"Pegleg Peet"
```


### Raffle with Custom Labels and CSV Output

#### CSV Input

##### `cat loot.csv`

```csv
"Reprobates of the Sea",Beard,Ham,Hook,Parrot,"Peg Leg",Soap
"Pegleg Peet",,,,,1,
"Captain Hook",,,2,1,,1
"The Pirate Captain",,3,,6,,
"The Surprisingly Curvaceous Pirate",,,,,,4
```

##### `csvtk pretty loot.csv`

```
Reprobates of the Sea                Beard   Ham   Hook   Parrot   Peg Leg   Soap
Pegleg Peet                                                        1         
Captain Hook                                       2      1                  1
The Pirate Captain                           3            6                  
The Surprisingly Curvaceous Pirate                                           4
```

#### Run Divvy

```
$ divvy -csv -player-label Pirates -loot-label Booty -token-label Doubloons loot.csv
Divvy Up the LOOT v1.0.0
Reading "loot.csv"
Progress: .....
"Pegleg Peet" got "Peg Leg"
"The Surprisingly Curvaceous Pirate" got "Soap"
"The Pirate Captain" got "Ham"
"Captain Hook" got "Hook"
CSV Output: DivvyUpTheLoot_2021-03-10_062633_PST.csv
```

#### CSV Output

##### `csvtk pretty DivvyUpTheLoot_2021-03-10_062633_PST.csv`

```
Pirates                              Booty     Booty List   Doubloons
Pegleg Peet                          Peg Leg   Ham          3
The Surprisingly Curvaceous Pirate   Soap      Hook         2
The Pirate Captain                   Ham       Parrot       7
Captain Hook                         Hook      Peg Leg      1
                                               Soap         5
```

