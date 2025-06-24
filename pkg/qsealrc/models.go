package qsealrc

import (
	"fmt"
	"os"
	"time"
)

const QsealrcFileName = "qsealrc.yaml"

type Qsealrc struct {
	Version             string   `yaml:"version" default:"1"`
	Namespace           string   `yaml:"namespace" required:"true"`
	ControllerName      string   `yaml:"controller_name" required:"true"`
	ControllerNamespace string   `yaml:"controller_namespace" required:"true"`
	Secrets             []Secret `yaml:"secrets"`
}

type Secret struct {
	// Name of the secret example: "my-secret"
	Name string `yaml:"name" required:"true"`
	// Path of the sealed secret example: "secrets/my-secret.env.sealed.yaml"
	Sealed *string `yaml:"sealed,omitempty"`
	// Path of the env secret example: "secrets/my-secret.env"
	Env *string `yaml:"env,omitempty"`
	// Path of the files secret example: "secrets/my-secret.yaml"
	Files []string `yaml:"files,omitempty"`
	// Type of the secret example: "kubernetes.io/dockerconfigjson
	Type string `yaml:"type" default:"Opaque" required:"true"`
}

func (s Secret) Validate() error {
	if s.Env != nil && len(s.Files) > 0 {
		return fmt.Errorf("you can only use one of env or files for qseal secret %s", s.Name)
	}
	return nil
}

func (secret Secret) SealedPath() (string, error) {
	if secret.Sealed != nil {
		return *secret.Sealed, nil
	}
	if secret.Env != nil {
		return *secret.Env + ".sealed.yaml", nil
	}
	if len(secret.Files) > 0 {
		return secret.Files[0] + ".sealed.yaml", nil
	}
	return "", fmt.Errorf("no sealed path could be determined for secret %s", secret.Name)
}

func (secret Secret) SyncFileTime() error {
	now := time.Now()
	sealedPath, err := secret.SealedPath()
	if err != nil {
		return err
	}
	err = os.Chtimes(sealedPath, now, now)
	if err != nil {
		return fmt.Errorf("error updating the date of the sealed secret %s: %v", sealedPath, err)
	}
	// we update the date of the env if not nil
	if secret.Env != nil {
		err = os.Chtimes(*secret.Env, now, now)
		if err != nil {
			return fmt.Errorf("error updating the date of the env %s: %v", *secret.Env, err)
		}
		return nil
	}
	// we update the date of the files
	for _, file := range secret.Files {
		err = os.Chtimes(file, now, now)
		if err != nil {
			return fmt.Errorf("error updating the date of the file %s: %v", file, err)
		}
	}
	return nil
}
