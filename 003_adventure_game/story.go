package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

type Story struct {
	arcs []StoryArc
}

func (s Story) printStory() {
	for _, arc := range s.arcs {
		fmt.Printf("Identifier: %v", arc.Identifier)
		fmt.Printf("Title: %v", arc.Title)
		fmt.Printf("Paragraph: %v", arc.Paragraph)
		for _, option := range arc.Options {
			fmt.Printf("Option Text: %v", option.Text)
			fmt.Printf("Option Arc: %v", option.Arc)
		}
	}
}

func (s *Story) Load(filePath string) error {
	jsonData, err := s.getJSON(filePath)

	if err != nil {
		log.Println("Unable to load JSON")
		return err
	}

	var arcs []StoryArc

	for key, data := range jsonData {
		arc := new(StoryArc)
		arc.Load(key, data.(map[string]any))
		arcs = append(arcs, *arc)
	}
	s.arcs = arcs

	return nil
}

func (s Story) getJSON(filePath string) (map[string]any, error) {
	file, err := os.Open(filePath)

	if err != nil {
		log.Printf("Cannot find %v file", filePath)
		return nil, err
	}

	defer file.Close()

	bytes, err := io.ReadAll(file)

	if err != nil {
		log.Printf("Unable to read %v file", filePath)
		return nil, err
	}

	var data any
	err = json.Unmarshal(bytes, &data)

	if err != nil {
		log.Printf("Unable to convert data into JSON")
		return nil, err
	}

	return data.(map[string]any), nil
}

func (s Story) GetArc(key string) (*StoryArc, error) {
	for _, arc := range s.arcs {
		if arc.Identifier == key {
			return &arc, nil
		}
	}
	return nil, errors.New("Cannot find " + key + " arc")
}

type StoryArc struct {
	Identifier string
	Title      string
	Paragraph  string
	Options    []*ArcOption
}

func (sa *StoryArc) Load(key string, data map[string]any) {
	var buffer bytes.Buffer
	paragraphs := data["story"].([]any)
	for _, p := range paragraphs {
		buffer.WriteString(p.(string) + "\r\n")
	}

	sa.Identifier = key
	sa.Title = data["title"].(string)
	sa.Paragraph = buffer.String()

	var options []*ArcOption

	for i, v := range data["options"].([]any) {
		opt := new(ArcOption)
		opt.Load(i+1, v.(map[string]interface{}))
		options = append(options, opt)
	}
	sa.Options = options
}

type ArcOption struct {
	Number int
	Text   string
	Arc    string
}

func (ao *ArcOption) Load(number int, data map[string]any) {
	ao.Number = number
	ao.Text = data["text"].(string)
	ao.Arc = data["arc"].(string)
}
