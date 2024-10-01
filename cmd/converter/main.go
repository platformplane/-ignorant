package main

import (
	"errors"

	"github.com/platformplane/scanner/pkg/config"
	"github.com/platformplane/scanner/pkg/converter"
)

func main() {
	cfg, err := config.Parse(".")

	if err != nil {
		panic(err)
	}

	c := converter.New(cfg, ".")

	var result error

	c.DeleteTrivyIgnore()

	if err := c.WriteTrivyIgnore(); err != nil {
		result = errors.Join(result, err)
	}

	if result != nil {
		panic(result)
	}
}
