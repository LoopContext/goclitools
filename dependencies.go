package goclitools

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/urfave/cli"
)

// DependencyScript interface - has to implement Run
type DependencyScript interface {
	Run() error
}

// DependencyScriptFn interface - has to implement Run
type DependencyScriptFn struct {
	Fn func() error
}

// Run runs the interface's fn
func (script DependencyScriptFn) Run() error {
	return script.Fn()
}

// DependencyScriptString interface
type DependencyScriptString struct {
	Fn string
}

// Run runs the dependency
func (script DependencyScriptString) Run() error {
	return RunInteractive(script.Fn)
}

// Dependency struct
type Dependency struct {
	Name               string
	CheckCmd           string
	CheckCmdValidation string
	Dependencies       []Dependency
	InstallScripts     []DependencyScript
	UninstallScripts   []DependencyScript
}

// Check checks the dependency
func (d *Dependency) Check() (bool, error) {

	output, err := Run(d.CheckCmd)

	if err != nil {
		if ee, ok := err.(*cli.ExitError); ok {
			if ee.ExitCode() != 1 {
				return false, ee
			}
		} else {
			return false, err
		}
	}

	if d.CheckCmdValidation != "" {
		if matched, err := regexp.Match(d.CheckCmdValidation, output); err != nil || !matched {
			return false, nil
		}
	}

	return string(output) != "", nil
}

// Install installs the dependency
func (d *Dependency) Install() error {
	res, _ := d.Check()
	if res == true {
		Logf("%s is already installed\n", d.Name)
		return nil
	}

	for _, dep := range d.Dependencies {
		Log("Validating subdependency", dep.Name)
		if err := dep.Install(); err != nil {
			return err
		}
	}

	count := len(d.InstallScripts)
	if count == 0 {
		return fmt.Errorf("%s cannot be installed (no install scripts)", d.Name)
	}

	for key, script := range d.InstallScripts {
		Logf("Running installation script %d/%d\n", key+1, count)
		if err := script.Run(); err != nil {
			return err
		}
	}

	Logf("Waiting for installation check to pass: ")
	attempts := 0
	for true {
		installed, err := d.Check()
		if err != nil {
			PrintNotOK()
			return err
		}
		if installed {
			break
		}

		if attempts > 60 {
			PrintNotOK()
			return errors.New("Installation check pass failed: timeout")
		}

		attempts++
		time.Sleep(time.Second)
	}
	PrintOK()

	return nil
}

// Uninstall uninstalls the dependency
func (d *Dependency) Uninstall() error {
	res, _ := d.Check()
	if res == false {
		Logf("%s is not installed\n", d.Name)
		return nil
	}
	count := len(d.UninstallScripts)
	if count == 0 {
		return fmt.Errorf("%s cannot be uninstalled (no uninstall scripts)", d.Name)
	}

	for key, script := range d.UninstallScripts {
		Logf("Running uninstallation script %d/%d\n", key+1, count)
		if err := script.Run(); err != nil {
			return err
		}
	}

	return nil
}
