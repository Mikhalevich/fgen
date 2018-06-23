package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	cDefaultSize = 4 * 1024
)

type Generator struct {
	params *Params
}

func NewGenerator(p *Params) *Generator {
	return &Generator{
		params: p,
	}
}

func (g *Generator) Start() {
	for _, fi := range g.params.Files {
		if fi.Count <= 0 {
			fi.Count = 1
		}

		for i := 0; i < fi.Count; i++ {
			path, err := g.makeFile(fi, i)
			if err != nil {
				fmt.Println(err)
				continue
			}

			err = g.changeTimes(path, fi)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (g *Generator) makeFile(fi FileInfo, number int) (string, error) {
	filePath := filepath.Join(g.params.Root, fmt.Sprintf("%s_%d", fi.Prefix, number))
	if fi.Suffix != "" {
		filePath = fmt.Sprintf("%s.%s", filePath, fi.Suffix)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if fi.Size < 0 {
		fi.Size = cDefaultSize
	}

	file.WriteString(strings.Repeat("test data", fi.Size/len("test data")))

	return filePath, nil
}

func (g *Generator) changeTimes(filePath string, fi FileInfo) error {
	aTime, err := g.parseTime(fi.ATime)
	if err != nil {
		return err
	}

	mTime, err := g.parseTime(fi.MTime)
	if err != nil {
		return err
	}

	err = os.Chtimes(filePath, aTime, mTime)
	if err != nil {
		return err
	}

	return nil
}

func (g *Generator) parseTime(t string) (time.Time, error) {
	if t == "" {
		return time.Now(), nil
	}

	return time.Parse(time.RFC3339, t)
}
