package main

import (
	"os"
	"testing"
)

const n = 10000

const inputFileName = "input.txt"
const outputFileName = "output.txt"

func BenchmarkStart(b *testing.B) {
	for i := 0; i < b.N; i++ {
		start(inputFileName, outputFileName, n)
	}

	os.Remove(outputFileName)
}
