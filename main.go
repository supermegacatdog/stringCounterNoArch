package main

import (
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 4 {
		log.Print(NotEnoughArgs)
		return
	}

	inputFileName := os.Args[1]
	outputFileName := os.Args[2]

	n, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Printf("failed to parse n:\n%s", err.Error())
		return
	}

	start(inputFileName, outputFileName, n)
}

func start(inputFileName string, outputFileName string, n int) {
	f, err := os.Open(inputFileName)
	if err != nil {
		log.Printf("failed to open input file:\n%s", err.Error())
		return
	}

	defer func() {
		er := f.Close()
		if err != nil {
			log.Printf("failed to close input file:\n%s", er.Error())
			return
		}
	}()

	inputFile := NewBigFile(f, n/2)

	for {
		block, er := inputFile.GetNextBlock()
		if er != nil {
			if er != io.EOF {
				log.Printf("failed to get next block:\n%s", er.Error())
				return
			}

			if len(block) == 0 && er == io.EOF {
				break
			}
		}

		er = inputFile.WriteBlock(block)
		if er != nil {
			log.Printf("failed to write block:\n%s", er.Error())
			return
		}
	}

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		log.Printf("failed to create output file:\n%s", err.Error())
		return
	}

	err = inputFile.MergeChunksAndSave(outputFile)
	if err != nil {
		log.Printf("failed to merge and save chunk files:\n%s", err.Error())
		return
	}

	err = outputFile.Close()
	if err != nil {
		log.Printf("failed to close output file:\n%s", err.Error())
		return
	}
}
