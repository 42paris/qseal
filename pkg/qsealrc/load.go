package qsealrc

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/mcuadros/go-defaults"
	"gopkg.in/yaml.v3"
)

func Load() (*Qsealrc, error) {
	_, err := os.Stat(QsealrcFileName)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file %s does not exist", QsealrcFileName)
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