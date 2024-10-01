package ignorant

import (
	"bytes"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func (c *Converter) DeleteTrivyIgnore() error {
	c.deleteFile(".trivyignore")
	c.deleteFile(".trivyignore.yml")
	c.deleteFile(".trivyignore.yaml")

	return nil
}

func (c *Converter) WriteTrivyIgnore() error {
	if err := c.writeTrivyIgnore(); err != nil {
		return err
	}

	if err := c.writeTrivyIgnoreYAML(); err != nil {
		return err
	}

	return nil
}

func (c *Converter) writeTrivyIgnore() error {
	b := new(bytes.Buffer)

	type Line struct {
		ID        string `yaml:"id,omitempty"`
		Statement string `yaml:"Statement,omitempty"`
	}

	var lines []Line

	for _, f := range c.Vulnerabilities {
		lines = append(lines, Line{
			ID:        f.ID,
			Statement: f.Statement,
		})
	}

	for _, f := range c.Misconfigurations {
		lines = append(lines, Line{
			ID:        f.ID,
			Statement: f.Statement,
		})
	}

	for i, l := range lines {
		if i > 0 {
			b.WriteString("\n")
		}

		if l.Statement != "" {
			b.WriteString("# ")
			b.WriteString(strings.TrimSpace(l.Statement))
			b.WriteString("\n")
		}

		b.WriteString(l.ID)
		b.WriteString("\n")
	}

	return c.writeFile(".trivyignore", b.Bytes())
}

func (c *Converter) writeTrivyIgnoreYAML() error {
	file := TrivyConfig{}

	for _, v := range c.Vulnerabilities {
		if v.ID == "" {
			continue
		}

		f := TrivyFinding{
			ID:        v.ID,
			Paths:     v.Paths,
			Statement: v.Statement,
		}

		if v.ExpiredAt != nil {
			f.ExpiredAt = &TrivyDate{Time: v.ExpiredAt}
		}

		file.Vulnerabilities = append(file.Vulnerabilities, f)
	}

	for _, m := range c.Misconfigurations {
		if m.ID == "" {
			continue
		}

		f := TrivyFinding{
			ID:        m.ID,
			Paths:     m.Paths,
			Statement: m.Statement,
		}

		if m.ExpiredAt != nil {
			f.ExpiredAt = &TrivyDate{Time: m.ExpiredAt}
		}

		file.Misconfigurations = append(file.Misconfigurations, f)
	}

	b := new(bytes.Buffer)

	encoder := yaml.NewEncoder(b)
	encoder.SetIndent(2)

	if err := encoder.Encode(file); err != nil {
		return err
	}

	return c.writeFile(".trivyignore.yaml", b.Bytes())
}

// https://github.com/aquasecurity/trivy/blob/7602d14654c2a4cca445439da7b93bfab4b6d872/pkg/result/ignore.go#L142
type TrivyConfig struct {
	Vulnerabilities   []TrivyFinding `yaml:"vulnerabilities,omitempty"`
	Misconfigurations []TrivyFinding `yaml:"misconfigurations,omitempty"`
	Secrets           []TrivyFinding `yaml:"secrets,omitempty"`
	Licenses          []TrivyFinding `yaml:"licenses,omitempty"`
}

// https://github.com/aquasecurity/trivy/blob/7602d14654c2a4cca445439da7b93bfab4b6d872/pkg/result/ignore.go#L24
type TrivyFinding struct {
	ID        string     `yaml:"id,omitempty"`
	Paths     []string   `yaml:"paths,omitempty"`
	ExpiredAt *TrivyDate `yaml:"expired_at,omitempty"`
	Statement string     `yaml:"statement,omitempty"`
}

type TrivyDate struct {
	*time.Time
}

func (d TrivyDate) MarshalYAML() (interface{}, error) {
	if d.Time == nil {
		return nil, nil
	}

	return yaml.Node{
		Kind:  yaml.ScalarNode,
		Style: yaml.FlowStyle,
		Value: d.Time.Format(time.DateOnly),
	}, nil
}
