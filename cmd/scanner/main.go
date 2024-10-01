package main

import (
	"errors"
	"os"
	"os/exec"

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
	defer c.DeleteTrivyIgnore()

	if err := c.WriteTrivyIgnore(); err != nil {
		result = errors.Join(result, err)
	}

	if result != nil {
		panic(result)
	}

	println("#")
	println("#")
	println("# VULNERABILITIES")
	println("#")

	vulnerabilities := exec.Command("trivy", "fs", "--scanners", "vuln", "--quiet", ".")
	vulnerabilities.Stdout = os.Stdout
	vulnerabilities.Stderr = os.Stderr

	vulnerabilities.Run()

	println("#")
	println("#")
	println("# MISCONFIGURATIONS")
	println("#")

	misconfigurations := exec.Command("trivy", "fs", "--scanners", "misconfig", "--quiet", ".")
	misconfigurations.Stdout = os.Stdout
	misconfigurations.Stderr = os.Stderr

	misconfigurations.Run()

	println("#")
	println("#")
	println("# SECRETS")
	println("#")

	gitleaks := exec.Command("gitleaks", "dir", "--no-banner", "--verbose")
	gitleaks.Stdout = os.Stdout
	gitleaks.Stderr = os.Stderr

	gitleaks.Run()

	println("#")
	println("#")
	println("# LICENSES")
	println("#")

	licenses := exec.Command("trivy", "fs", "--scanners", "license", "--quiet", ".")
	licenses.Stdout = os.Stdout
	licenses.Stderr = os.Stderr

	licenses.Run()
}
