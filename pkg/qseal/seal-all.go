package qseal

import (
	"fmt"
	"os"

	"github.com/42paris/qseal/pkg/qsealrc"
)

func SealAll(qsealRc qsealrc.Qsealrc) error {
	sealClient, err := NewKubeSealClient(qsealRc)
	if err != nil {
		return err
	}
	sealedPaths := make(map[string]bool)
	for _, secret := range qsealRc.Secrets {
		sealedPath, err := secret.SealedPath()
		if err != nil {
			return err
		}
		// we check if the path is already in the map
		_, ok := sealedPaths[sealedPath]
		if !ok {
			// we clear the file
			err = os.WriteFile(sealedPath, []byte{}, 0644)
			if err != nil {
				return fmt.Errorf("error clearing file %s: %v", sealedPath, err)
			}
			sealedPaths[sealedPath] = true
		}

		err = sealClient.Seal(secret)
		if err != nil {
			return err
		}
	}
	return nil
}
