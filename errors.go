package main

import "errors"

var (
	NotEnoughArgs   = errors.New("not enough args")
	StringIncorrect = errors.New("string is incorrect")
)
