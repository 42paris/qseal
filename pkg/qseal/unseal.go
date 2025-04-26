package qseal

import (
	"crypto/rsa"
	"fmt"
	"io"
	"os"

	"github.com/42paris/qseal/pkg/qsealrc"
	"github.com/42paris/qseal/pkg/secretgen"
	ssv1alpha1 "github.com/bitnami-labs/sealed-secrets/pkg/apis/sealedsecrets/v1alpha1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
)

func (*KubeSealClient) Unseal(secret qsealrc.Secret, keySet map[string]*rsa.PrivateKey) error {
	// Load the sealed secret
	sealedPath, err := secret.SealedPath()
	if err != nil {
		return err
	}
	r, err := os.Open(sealedPath)
	if err != nil {
		return err
	}

	// Decode all the sealed secrets
	sealedSecrets, err := decodeSealedSecrets(r)
	if err != nil {
		return err
	}
	// we find the sealedSecret with the name
	// this could be optimized with a cache but this should be done in the
	// parent function
	var sealedSecret *ssv1alpha1.SealedSecret
	for _, ss := range sealedSecrets {
		if ss.Name == secret.Name {
			sealedSecret = ss
			break
		}
	}
	if sealedSecret == nil {
		return fmt.Errorf("sealed secret %s not found for %s", secret.Name, sealedPath)
	}
	unsealed, err := sealedSecret.Unseal(scheme.Codecs, keySet)
	if err != nil {
		return err
	}
	err = secretgen.DeGen(unsealed, secret)
	if err != nil {
		return err
	}
	// we update the date of the files
	// this allow us to sync smartly the sealed secret with the source and vice versa
	err = secret.SyncFileTime()
	if err != nil {
		return fmt.Errorf("error updating the date of the secrets %s: %v", sealedPath, err)
	}
	return nil
}

func decodeSealedSecrets(r io.Reader) ([]*ssv1alpha1.SealedSecret, error) {
	decoder := yaml.NewYAMLOrJSONDecoder(r, 4096)
	sealedSecrets := []*ssv1alpha1.SealedSecret{}
	for {
		ss := &ssv1alpha1.SealedSecret{}
		err := decoder.Decode(ss)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		sealedSecrets = append(sealedSecrets, ss)
	}
	return sealedSecrets, nil
}
