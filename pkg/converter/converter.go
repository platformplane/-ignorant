package converter

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/platformplane/scanner/pkg/config"
)

type Converter struct {
	*config.Config

	root string
}

func New(root string) (*Converter, error) {
	cfg, err := config.Parse(root)

	if err != nil {
		return nil, err
	}

	return &Converter{
		cfg,
		root,
	}, nil
}

func (c *Converter) EnsureIngoreFiles() error {
	c.DeleteIngoreFiles()

	var result error

	if err := c.WriteTrivyFiles(); err != nil {
		result = errors.Join(result, err)
	}

	return result
}

func (c *Converter) DeleteIngoreFiles() error {
	c.DeleteTrivyFiles()

	return nil
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
