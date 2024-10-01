package main

import (
	"errors"

	"github.com/platformplane/ignorant"
)

func main() {
	cfg, err := ignorant.Parse(".")

	if err != nil {
		panic(err)
	}

	c := ignorant.New(cfg, ".")

	var result error

	if err := c.WriteTrivyIgnore(); err != nil {
		result = errors.Join(result, err)
	}

	if result != nil {
		panic(result)
	}
}
