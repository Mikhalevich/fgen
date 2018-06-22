package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/Mikhalevich/argparser"
)

type Params struct {
	Root  string     `json:"root, omitempty"`
	Files []FileInfo `json:"files"`
}

type FileInfo struct {
	Prefix string `json:"prefix"`
	ATime  string `json:"atime, omitempty"`
	MTime  string `json:"mtime, omitempty"`
	Size   int    `json:"size"`
	Count  int    `json:"count"`
}

func NewParams() *Params {
	return &Params{
		Root: ".",
		Files: []FileInfo{FileInfo{
			Prefix: "testFile",
			ATime:  time.Now().Format(time.RFC3339),
			MTime:  time.Now().Format(time.RFC3339),
			Size:   1024 * 1024,
			Count:  100,
		},
		},
	}
}

func loadParams() (*Params, error) {
	basicParams := NewParams()
	p, err, gen := argparser.Parse(basicParams)

	if gen {
		return nil, errors.New("Config was generated")
	}

	params := p.(*Params)

	if argparser.NArg() > 0 {
		params.Root = argparser.Arg(0)
	}

	if params.Root == "" {
		return nil, errors.New("Root path is not specified")
	}

	return params, err
}

func main() {
	startTime := time.Now()

	params, err := loadParams()
	if err != nil {
		fmt.Println(err)
		return
	}

	g := NewGenerator(params)
	g.Start()

	fmt.Printf("Execution time = %v\n", time.Now().Sub(startTime))
}
