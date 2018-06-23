package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/Mikhalevich/argparser"
	"gopkg.in/cheggaaa/pb.v1"
)

type Params struct {
	Root  string     `json:"root, omitempty"`
	Files []FileInfo `json:"files"`
}

type FileInfo struct {
	Prefix string `json:"prefix"`
	Suffix string `json:"suffix, omitempty"`
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
			Suffix: "tst",
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

func showProgress(c chan int64) {
	go func() {
		content := <-c

		bar := pb.New64(content)
		bar.Start()

		for chunk := range c {
			bar.Add64(chunk)
		}
		bar.Finish()
	}()
}

func main() {
	startTime := time.Now()

	params, err := loadParams()
	if err != nil {
		fmt.Println(err)
		return
	}

	g := NewGenerator(params)
	showProgress(g.Notifier)
	errs := g.Start()

	if len(errs) > 0 {
		fmt.Println("Errors:")
		for _, err := range errs {
			fmt.Printf("Error: %v\n", err)
		}
	}

	fmt.Printf("Execution time = %v\n", time.Now().Sub(startTime))
}
