package qsealrc

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/mcuadros/go-defaults"
	"gopkg.in/yaml.v3"
)

func walkDir() (dir string, found bool, err error) {
	dir, err = os.Getwd()
	if err != nil {
		return "", false, fmt.Errorf("error reading working directory")
	}

	for {
		rcFullpath := filepath.Join(dir, QsealrcFileName)
		_, err := os.Stat(rcFullpath)
		if err == nil {
			found = true
			break
		} else {
			if !os.IsNotExist(err) {
				return "", false, fmt.Errorf("error searching for %s: %w", QsealrcFileName, err)
			}

			nextDir := filepath.Dir(dir)
			if nextDir == dir { // Stop at root dir
				break
			}
			dir = nextDir
		}
	}

	return dir, found, nil
}

func Load() (*Qsealrc, error) {
	dir, found, err := walkDir()
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("file %s does not exist", QsealrcFileName)
	}

	err = os.Chdir(dir)
	if err != nil {
		return nil, fmt.Errorf("could not change directory to %s: %w", dir, err)
	}

	file, err := os.Open(QsealrcFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var qsealrc Qsealrc
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&qsealrc)
	if err != nil {
		return nil, err
	}
	// Set default values
	defaults.SetDefaults(&qsealrc)
	// Validate the struct
	validate := validator.New()
	err = validate.Struct(qsealrc)
	if err != nil {
		return nil, err
	}
	// Validate the secrets
	for _, secret := range qsealrc.Secrets {
		err = secret.Validate()
		if err != nil {
			return nil, err
		}
	}
	return &qsealrc, nil	
}