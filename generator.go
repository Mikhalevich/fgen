package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Mikhalevich/jober"
)

const (
	cDefaultSize    = 4 * 1024
	cDefaulfWorkers = 100
)

type Generator struct {
	params   *Params
	Notifier chan int64
}

func NewGenerator(p *Params) *Generator {
	if p.Workers <= 0 {
		p.Workers = cDefaulfWorkers
	}
	return &Generator{
		params:   p,
		Notifier: make(chan int64, p.Workers),
	}
}

func (g *Generator) Start() []error {
	var totalFiles int64 = 0
	for _, fi := range g.params.Files {
		if fi.Count <= 0 {
			fi.Count = 1
		}
		totalFiles += int64(fi.Count)
	}
	g.Notifier <- totalFiles

	fileJob := jober.NewWorkerPool(jober.NewAll(), g.params.Workers)

	for _, fi := range g.params.Files {
		for i := 0; i < fi.Count; i++ {
			index := i
			workerFunc := func() (interface{}, error) {
				path, err := g.makeFile(fi, index)
				if err != nil {
					return nil, err
				}

				err = g.changeTimes(path, fi)
				if err != nil {
					return nil, err
				}

				g.Notifier <- 1

				return nil, nil
			}
			fileJob.Add(workerFunc)
		}
	}

	fileJob.Wait()

	_, errs := fileJob.Get()
	return errs
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
