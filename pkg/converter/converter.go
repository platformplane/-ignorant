package converter

import (
	"os"
	"path/filepath"

	"github.com/platformplane/scanner/pkg/config"
)

type Converter struct {
	*config.Config

	root string
}

func New(cfg *config.Config, root string) *Converter {
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

func (c *Converter) deleteFile(name string) error {
	root, err := filepath.Abs(c.root)

	if err != nil {
		return err
	}

	path := filepath.Join(root, name)

	os.Remove(path)

	return nil
}
