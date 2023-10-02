package srt

import (
	"bufio"
	"os"
	"strings"
)

type Subtitle struct {
	Index string `json:"index"`
	Start string `json:"start"`
	End   string `json:"end"`
	Text  string `json:"text"`
}

func SrtToJson(srtPath string) ([]Subtitle, error) {
	file, err := os.Open(srtPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var subs []Subtitle
	var sub Subtitle
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			subs = append(subs, sub)
			sub = Subtitle{}
		} else if sub.Index == "" {
			sub.Index = line
		} else if sub.Start == "" {
			times := strings.Split(line, " --> ")
			sub.Start = times[0]
			sub.End = times[1]
		} else {
			sub.Text += line + "\n"
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return subs, nil
}
