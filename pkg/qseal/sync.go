package qseal

import (
	"fmt"
	"os"

	"gitlab.42paris.fr/froz/qseal/pkg/qsealrc"
)

func Sync(qsealRc qsealrc.Qsealrc) error {
	secretsBySealedPath, actionBySealedPath, err := getSecretsStatus(qsealRc)
	if err != nil {
		return fmt.Errorf("error getting secrets status: %v", err)
	}
	sealClient, err := NewKubeSealClient(qsealRc)
	if err != nil {
		return fmt.Errorf("error creating seal client: %v", err)
	}
	keySet, err := sealClient.RetrievePrivateKeys()
	if err != nil {
		return fmt.Errorf("error retrieving private keys: %v", err)
	}
	for sealedPath, secrets := range secretsBySealedPath {
		logSecretAction(actionBySealedPath[sealedPath], sealedPath, secrets)
		switch actionBySealedPath[sealedPath] {
		case SyncActionSeal:
			err = os.WriteFile(sealedPath, []byte{}, 0644)
			if err != nil {
				return fmt.Errorf("error clearing file %s: %v", sealedPath, err)
			}
			for _, secret := range secrets {
				err = sealClient.Seal(secret)
				if err != nil {
					return fmt.Errorf("error sealing secret %s: %v", secret.Name, err)
				}
			}
		case SyncActionUnseal:
			for _, secret := range secrets {
				err = sealClient.Unseal(secret, keySet)
				if err != nil {
					return fmt.Errorf("error unsealing secret %s: %v", secret.Name, err)
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
