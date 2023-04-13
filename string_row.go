package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Row struct {
	Query     string
	Frequency int64
}

func (r *Row) Parse(str string) error {
	parsedString := strings.Split(str, " ")
	if len(parsedString) < 2 {
		return StringIncorrect
	}

	freq, err := strconv.Atoi(parsedString[1])
	if err != nil {
		return fmt.Errorf("failed to parse frequency: %w\n", err)
	}

	r.Query = parsedString[0]
	r.Frequency = int64(freq)

	return nil
}
