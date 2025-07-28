package main

import (
	"fmt"
	"log"
	"os"
)

type ConsoleRunner struct {
}

func (cr ConsoleRunner) Start(provider *StoryArcProvider) {
	cr.displayArcText(*provider, "intro")
}

func (cr ConsoleRunner) displayArcText(provider StoryArcProvider, arcName string) {

	arc, err := provider.WriteTemplatedText(os.Stdout, arcName)
	if err != nil {
		log.Println(err)
	}
	if len(arc.Options) == 0 {
		return
	}
	fmt.Print("Your Option: ")
	var optionNumber int
	fmt.Scan(&optionNumber)
	for _, option := range arc.Options {
		if option.Number == optionNumber {
			cr.displayArcText(provider, option.Arc)
		}
	}
}
