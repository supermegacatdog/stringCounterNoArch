package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type ChunkFile struct {
	File    *os.File
	Scanner *bufio.Scanner
	Max     int
}

func NewChunkFile(file *os.File, max int) *ChunkFile {
	return &ChunkFile{
		File:    file,
		Scanner: bufio.NewScanner(file),
		Max:     max,
	}
}

func (f *ChunkFile) Rewrite(group Group) error {
	stat, err := f.File.Stat()
	if err != nil {
		return fmt.Errorf("failed to get temp file stat: %w\n", err)
	}

	// if a temp file is new, its is much more convenient to save data here and return
	if stat.Size() == 0 {
		_, er := f.File.Write(group.BytesEncode())
		if er != nil {
			return fmt.Errorf("failed to write to new file: %w\n", er)
		}

		group.Clear()
		return nil
	}

	content := make([]byte, stat.Size())
	_, err = f.File.ReadAt(content, 0)
	if err != nil {
		return fmt.Errorf("failed to get read file: %w\n", err)
	}

	tempQueries := make(Group, 0)

	for _, row := range strings.Split(strings.Trim(string(content), "\n"), "\n") {
		if row == "" {
			continue
		}
		parsedString := &Row{}
		er := parsedString.Parse(row)
		if er != nil {
			return fmt.Errorf("failed to parse string of temp file: %w\n", er)
		}

		tempQueries[parsedString.Query] = parsedString.Frequency
	}

	tempQueries.Add(group, f.Max)

	_, er := f.File.WriteAt(tempQueries.BytesEncode(), 0)
	if er != nil {
		return fmt.Errorf("failed to write temp file at: %w\n", er)
	}

	return nil
}
