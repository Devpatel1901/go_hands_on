package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"strings"
	"time"
)

func readCSVFile(filename string) [][]string {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal("Error while reading the file: ", err)
	}

	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()

	if err != nil {
		fmt.Println("Error reading records")
	}

	return records
}

func readAnswer() string {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}

	input = strings.TrimSpace(strings.ToLower(input))

	return input
}

func main() {

	fileName := flag.String("f", "problems.csv", "a csv file in the format of 'question, answer' (default 'problems.csv')")
	timeout := flag.Int("t", 30, "the time limit for quiz in seconds (default 30)")
	isShuffle := flag.Bool("s", false, "shuffle the order of questions or not (default false)")

	flag.Parse()

	if *timeout < 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	records := readCSVFile(*fileName)

	if *isShuffle {
		rand.Shuffle(len(records), func(i, j int) {
			records[i], records[j] = records[j], records[i]
		})
	}

	doneCh := make(chan bool)
	total := 0

	// Go Routine 1 --> For Q&A
	go func() {
		for idx, record := range records {
			fmt.Printf("Problem #%d: %v = ", idx, record[0])

			answer := readAnswer()
			expectedAnswer := record[1]

			if answer == expectedAnswer {
				total++
			}
		}
		doneCh <- true
	}()

	// Go Routine 2 --> For time keeping
	go func() {
		timer := time.NewTimer(time.Duration(*timeout) * time.Second)
		<-timer.C
		doneCh <- true
	}()

	<-doneCh
	fmt.Println("")
	fmt.Printf("You scored %d out of %d", total, len(records))
}
