package secretgen

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/42paris/qseal/pkg/qsealrc"
	corev1 "k8s.io/api/core/v1"
)

func DeGen(unsealed *corev1.Secret, secret qsealrc.Secret) error {
	if secret.Env != nil {
		envPath := *secret.Env
		file, err := os.OpenFile(envPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("failed to open env file %s: %w", envPath, err)
		}
		defer file.Close()
		writer := io.Writer(file)
		for key, value := range unsealed.Data {
			// Write the key-value pair to the env file
			_, err := fmt.Fprintf(writer, "%s=%s\n", key, value)
			if err != nil {
				return fmt.Errorf("failed to write to env file %s: %w", envPath, err)
			}
		}
		return nil
	}
	// in the case of files
	for _, filePath := range secret.Files {
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", filePath, err)
		}
		defer file.Close()
		// we keep the file name without the path
		fileName := filepath.Base(filePath)
		data, ok := unsealed.Data[fileName]
		if !ok {
			return fmt.Errorf("file %s not found in secret data for %s", fileName, secret.Name)
		}
		_, err = file.Write(data)
		if err != nil {
			return fmt.Errorf("failed to write to file %s: %w", filePath, err)
		}
	}
	return nil
}
