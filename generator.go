package main

import (
	"fmt"
	"math/rand"
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
			g.makeFile(fi, i)
		}
	}
}

func (g *Generator) makeFile(fi FileInfo, number int) {
	filePath := filepath.Join(g.params.Root, fmt.Sprintf("%s_%d", fi.Prefix, number))
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	if fi.MaxSize < fi.MinSize {
		fi.MaxSize = fi.MinSize
	}

	sizeDiff := fi.MaxSize - fi.MinSize
	delta := 0
	if sizeDiff > 0 {
		delta = rand.Intn(sizeDiff)
	}
	size := fi.MinSize + delta

	file.WriteString(strings.Repeat("test data", size/len("test data")))

	var aTime time.Time
	if fi.ATime != "" {
		aTime, err = time.Parse(time.RFC3339, fi.ATime)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		aTime = time.Now()
	}

	var mTime time.Time
	if fi.MTime != "" {
		mTime, err = time.Parse(time.RFC3339, fi.MTime)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		mTime = time.Now()
	}

	err = os.Chtimes(filePath, aTime, mTime)
	if err != nil {
		fmt.Println(err)
	}
}
