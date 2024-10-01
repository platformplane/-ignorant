package ignorant

import (
	"os"
	"path/filepath"
)

type Converter struct {
	*Config

	root string
}

func New(cfg *Config, root string) *Converter {
	return &Converter{
		cfg,
		root,
	}
}

func (c *Converter) writeFile(name string, data []byte) error {
	root, err := filepath.Abs(c.root)

	if err != nil {
		return err
	}

	path := filepath.Join(root, name)

	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}

	return nil

}
