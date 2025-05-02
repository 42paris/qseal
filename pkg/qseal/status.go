package qseal

import (
	"fmt"
	"os"
	"time"

	"github.com/42paris/qseal/pkg/qsealrc"
	"github.com/fatih/color"
)

var (
	warnLog   = color.New(color.FgYellow).SprintFunc()
	actionLog = color.New(color.FgGreen).SprintFunc()
	pathLog   = color.New(color.FgHiBlue).SprintFunc()
)

type SyncAction string

const (
	SyncActionUnseal    SyncAction = "UNSEALING"
	SyncActionSeal      SyncAction = "SEALING"
	SyncActionDoNothing SyncAction = "SKIP"
)

const timeFormat = "2006-01-02 15:04:05"

func Status(qsealRc qsealrc.Qsealrc) error {
	secretsBySealedPath, actionBySealedPath, err := getSecretsStatus(qsealRc)
	if err != nil {
		return fmt.Errorf("error getting secrets status: %v", err)
	}
	for sealedPath, secrets := range secretsBySealedPath {
		logSecretAction(actionBySealedPath[sealedPath], sealedPath, secrets)
	}
	return nil
}

func logSecretAction(action SyncAction, sealedPath string, secrets []qsealrc.Secret) {
	actionLogFunc := actionLog
	if action == SyncActionDoNothing {
		actionLogFunc = warnLog
	}
	fmt.Printf("[%s] %s %s (%d secret(s))\n",
		time.Now().Format(timeFormat),
		actionLogFunc(action),
		pathLog(sealedPath),
		len(secrets),
	)
}

func getSecretsStatus(qsealRc qsealrc.Qsealrc) (map[string][]qsealrc.Secret, map[string]SyncAction, error) {
	secretsBySealedPath := make(map[string][]qsealrc.Secret)
	actionBySealedPath := make(map[string]SyncAction)
	for _, secret := range qsealRc.Secrets {
		sealedPath, err := secret.SealedPath()
		if err != nil {
			return nil, nil, fmt.Errorf("error getting sealed path for secret %s: %v", secret.Name, err)
		}
		action, err := decideSyncAction(sealedPath, secret)
		if err != nil {
			return nil, nil, fmt.Errorf("error deciding sync action for secret %s: %v", secret.Name, err)
		}
		existingAction, ok := actionBySealedPath[sealedPath]
		if ok && action == SyncActionDoNothing {
			// if we have already a action for this sealed path and
			// the action is do nothing we don't overwrite it
			continue
		}

		// is have already a action for this sealed path that is not
		// the same as the one we are trying to add this mean we have a conflict
		if ok && existingAction != action {
			return nil, nil, fmt.Errorf(
				"conflict for sealed secret %s: %s (%s) and %s (%s) try resolve the conflict by seal-all or unseal-all",
				sealedPath,
				existingAction,
				secret.Name, action,
				secret.Name)
		}
		actionBySealedPath[sealedPath] = action
		secretsBySealedPath[sealedPath] = append(secretsBySealedPath[sealedPath], secret)
	}
	return secretsBySealedPath, actionBySealedPath, nil
}

func decideSyncAction(sealedPath string, secret qsealrc.Secret) (SyncAction, error) {
	// we check if the sealed secret file exists
	// if it does not exist, we return SyncActionSeal
	sealedFileDate, err := getFileDate(sealedPath)
	if err != nil {
		if os.IsNotExist(err) {
			return SyncActionSeal, nil
		}
		return SyncActionDoNothing, fmt.Errorf("error checking file %s: %v", sealedPath, err)
	}
	// if we have a env file we gatter the date of the file
	var sourceFileDate time.Time
	if secret.Env != nil {
		sourceFileDate, err = getFileDate(*secret.Env)
		if err != nil {
			if os.IsNotExist(err) {
				return SyncActionUnseal, nil
			}
			return SyncActionDoNothing, fmt.Errorf("error checking file %s: %v", *secret.Env, err)
		}
	} else {
		// we check the files and we keep the most recent date
		for _, file := range secret.Files {
			fileDate, err := getFileDate(file)
			if err != nil {
				if os.IsNotExist(err) {
					return SyncActionUnseal, nil
				}
				return SyncActionDoNothing, fmt.Errorf("error checking file %s: %v", file, err)
			}
			if fileDate.After(sourceFileDate) {
				sourceFileDate = fileDate
			}
		}
	}
	// if the source file is more recent than the sealed file, we return SyncActionSeal
	if sourceFileDate.After(sealedFileDate) {
		return SyncActionSeal, nil
	}
	// if the sealed file is more recent than the source file, we return SyncActionUnseal
	if sealedFileDate.After(sourceFileDate) {
		return SyncActionUnseal, nil
	}
	// if the files are the same date, we return SyncActionDoNothing
	return SyncActionDoNothing, nil
}

func getFileDate(path string) (time.Time, error) {
	fileStats, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	if fileStats.IsDir() {
		return time.Time{}, fmt.Errorf("error: %s is a directory", path)
	}
	return fileStats.ModTime(), nil
}
