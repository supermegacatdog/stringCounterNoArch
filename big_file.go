package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type BigFile struct {
	File       *os.File
	Scanner    *bufio.Scanner
	ChunkFiles []*ChunkFile
	max        int
	lastToken  string
}

func NewBigFile(file *os.File, max int) *BigFile {
	return &BigFile{
		File:    file,
		Scanner: bufio.NewScanner(file),
		max:     max,
	}
}

func (f *BigFile) GetNextBlock() (Group, error) {
	queries := make(Group, 0)

	if f.lastToken != "" {
		queries[f.lastToken]++
	}

	for f.Scanner.Scan() {
		// filling inQueries map using scan
		// map cannot include more than N/2 elements
		if len(queries) < f.max {
			queries[f.Scanner.Text()]++
			continue
		}

		f.lastToken = f.Scanner.Text()
		return queries, nil
	}

	f.lastToken = ""
	return queries, io.EOF
}

func (f *BigFile) WriteBlock(block Group) error {
	for i := 0; len(block) > 0; i++ {
		if i > len(f.ChunkFiles)-1 {
			tempFile, err := os.Create(fmt.Sprintf(TempFileNameMask, i))
			if err != nil {
				return fmt.Errorf("failed to create temp file:, %w\n", err)
			}

			f.ChunkFiles = append(f.ChunkFiles, NewChunkFile(tempFile, f.max))
		}

		err := f.ChunkFiles[i].Rewrite(block)
		if err != nil {
			return fmt.Errorf("failed to rewrite chunk file: %w\n", err)
		}
	}

	return nil
}

func (f *BigFile) MergeChunksAndSave(outputFile *os.File) error {
	for _, file := range f.ChunkFiles {
		_, err := file.File.Seek(0, 0)
		if err != nil {
			return fmt.Errorf("failed to get current position: %w\n", err)
		}

		_, err = io.Copy(outputFile, file.File)
		if err != nil {
			return fmt.Errorf("failed to copy file into another: %w\n", err)
		}

		err = os.Remove(file.File.Name())
		if err != nil {
			return fmt.Errorf("failed to remove file: %w\n", err)
		}
	}

	return nil
}
