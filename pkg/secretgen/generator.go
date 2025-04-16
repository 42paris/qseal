package secretgen

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gitlab.42paris.fr/froz/qseal/pkg/kubemodels"
	"gitlab.42paris.fr/froz/qseal/pkg/qsealrc"
)

func Gen(secret qsealrc.Secret) (*kubemodels.Secret, error) {
	kubeSecret := &kubemodels.Secret{
		APIVersion: "v1",
		Kind:       "Secret",
		Type:       secret.Type,
		Metadata: kubemodels.Metadata{
			Name: secret.Name,
		},
		Data: make(map[string]string),
	}

	// in the case of a env file
	if secret.Env != nil {
		envPath := *secret.Env
		file, err := os.Open(envPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open env file %s: %v", envPath, err)
		}
		defer file.Close()
		reader := io.Reader(file)
		env, err := parseEnvFile(reader)
		if err != nil {
			return nil, err
		}
		for key, value := range env {
			// Encode the value to base64
			encodedValue := base64.StdEncoding.EncodeToString([]byte(value))
			// Add the encoded value to the secret data
			kubeSecret.Data[key] = encodedValue
		}
		return kubeSecret, nil
	}
	
	// in the case of files
	for _, filePath := range secret.Files {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %v", filePath, err)
		}
		defer file.Close()
		// we keep the file name without the path
		fileName := filepath.Base(filePath)
		data, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}
		encodedValue := base64.StdEncoding.EncodeToString(data)
		kubeSecret.Data[fileName] = encodedValue
	}
	return kubeSecret, nil
}

func parseEnvFile(file io.Reader) (map[string]string, error) {
	env := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line[0] == '#' {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid env line: %s", line)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		env[key] = value
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return env, nil
}
