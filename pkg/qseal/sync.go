package qseal

import (
	"fmt"
	"os"

	"github.com/42paris/qseal/pkg/qsealrc"
)

func Sync(qsealRc qsealrc.Qsealrc) error {
	secretsBySealedPath, actionBySealedPath, err := getSecretsStatus(qsealRc)
	if err != nil {
		return fmt.Errorf("error getting secrets status: %w", err)
	}
	sealClient, err := NewKubeSealClient(qsealRc)
	if err != nil {
		return fmt.Errorf("error creating seal client: %w", err)
	}
	keySet, err := sealClient.RetrievePrivateKeys()
	if err != nil {
		return fmt.Errorf("error retrieving private keys: %w", err)
	}
	for sealedPath, secrets := range secretsBySealedPath {
		logSecretAction(actionBySealedPath[sealedPath], sealedPath, secrets)
		switch actionBySealedPath[sealedPath] {
		case SyncActionSeal:
			err = os.WriteFile(sealedPath, []byte{}, 0644)
			if err != nil {
				return fmt.Errorf("error clearing file %s: %w", sealedPath, err)
			}
			for _, secret := range secrets {
				err = sealClient.Seal(secret)
				if err != nil {
					return fmt.Errorf("error sealing secret %s: %w", secret.Name, err)
				}
			}
		case SyncActionUnseal:
			for _, secret := range secrets {
				err = sealClient.Unseal(secret, keySet)
				if err != nil {
					return fmt.Errorf("error unsealing secret %s: %w", secret.Name, err)
				}
			}
		case SyncActionDoNothing:
			// do nothing
		default:
			panic(fmt.Sprintf("unknown action %s for %s", actionBySealedPath[sealedPath], sealedPath))
		}
	}
	return nil
}
