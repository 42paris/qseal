package qseal

import (
	"fmt"
	"time"

	"github.com/42paris/qseal/pkg/qsealrc"
)

func UnsealAll(qsealRc qsealrc.Qsealrc) error {
	sealClient, err := NewKubeSealClient(qsealRc)
	if err != nil {
		return err
	}
	keySet, err := sealClient.RetrievePrivateKeys()
	if err != nil {
		return err
	}
	now := time.Now()
	// we dump a json the keySet
	for _, secret := range qsealRc.Secrets {
		err := sealClient.Unseal(secret, now, keySet)
		if err != nil {
			return fmt.Errorf("error unsealing secret %s: %v", secret.Name, err)
		}
	}
	return nil
}
