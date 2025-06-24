package qseal

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"

	"github.com/42paris/qseal/pkg/qsealrc"
	"github.com/42paris/qseal/pkg/secretgen"
	"github.com/bitnami-labs/sealed-secrets/pkg/apis/sealedsecrets/v1alpha1"
	"github.com/bitnami-labs/sealed-secrets/pkg/kubeseal"
	"gopkg.in/yaml.v3"
	"k8s.io/client-go/kubernetes/scheme"
)

func (k *KubeSealClient) Seal(secret qsealrc.Secret) error {
	err := secret.Validate()
	if err != nil {
		return err
	}
	sealedPath, err := secret.SealedPath()
	if err != nil {
		return err
	}
	out, err := os.OpenFile(sealedPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, fs.ModePerm)
	if err != nil {
		return err
	}
	defer out.Close()

	// First we transform the file into a kubernetes secret
	kubeSecret, err := secretgen.Gen(secret)
	if err != nil {
		return err
	}

	// Then we encode the kubernetes secret into a yaml file
	rawSecretBuffer := &bytes.Buffer{}
	yamlEncoder := yaml.NewEncoder(rawSecretBuffer)
	yamlEncoder.SetIndent(2)
	err = yamlEncoder.Encode(kubeSecret)
	if err != nil {
		return err
	}

	// Then we we seal the secret

	err = kubeseal.Seal(
		k.clientConfig,
		"yaml",
		rawSecretBuffer,
		out,
		scheme.Codecs,
		k.pubKey, // the pubKey
		v1alpha1.StrictScope,
		false, // allow empty data
		secret.Name,
		k.qsealrc.Namespace,
	)
	if err != nil {
		return err
	}
	// we update the date of the files
	err = secret.SyncFileTime()
	if err != nil {
		return fmt.Errorf("error updating the date of the secrets %s: %w", sealedPath, err)
	}
	return nil
}
