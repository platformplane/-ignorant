package main

import (
	"github.com/platformplane/scanner/pkg/converter"
)

func main() {
	c, err := converter.New(".")

	if err != nil {
		panic(err)
	}

	if err := c.EnsureIngoreFiles(); err != nil {
		panic(err)
	}
}
