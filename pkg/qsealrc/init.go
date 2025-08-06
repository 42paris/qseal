package qsealrc

import (
	_ "embed"
	"fmt"
	"os"
	"text/template"
	"time"
)

//go:embed header.yaml
var header string

type HeaderTemplate struct {
	Author string
	Date   string
}

func Init(ignoreParents bool) error {
	// Check if the .qsealrc.yaml file exists in current directory
	_, err := os.Stat(QsealrcFileName)
	if !os.IsNotExist(err) {
		return fmt.Errorf("file %s already exists in current directory", QsealrcFileName)
	}
	// Check if the .qsealrc.yaml file exists in parent directories (unless ignored)
	if !ignoreParents {
		dir, found, err := walkDir()
		if err != nil {
			return err
		}
		if found {
			return fmt.Errorf("file %s already exists in parent directory %s (use -i to ignore parents)", QsealrcFileName, dir)
		}
	}
	// Create the .qsealrc.yaml file with default values
	file, err := os.Create(QsealrcFileName)
	if err != nil {
		return err
	}
	defer file.Close()
	// Parse the template
	tmpl, err := template.New("header").Parse(header)
	if err != nil {
		return err
	}
	headerData := HeaderTemplate{
		Author: os.Getenv("USER"),
		Date:   time.Now().In(time.Local).Format("2006-01-02 15:04:05"),
	}
	return tmpl.Execute(file, headerData)
}
