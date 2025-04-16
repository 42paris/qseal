package qseal

import (
	"context"
	"crypto/rsa"
	"os"

	"github.com/bitnami-labs/sealed-secrets/pkg/kubeseal"
	"gitlab.42paris.fr/froz/qseal/pkg/qsealrc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeSealClient struct {
	clientConfig clientcmd.ClientConfig
	qsealrc      qsealrc.Qsealrc
	pubKey       *rsa.PublicKey
}

func NewKubeSealClient(
	qsealrc qsealrc.Qsealrc,
) (*KubeSealClient, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	clientConfig := clientcmd.NewInteractiveDeferredLoadingClientConfig(
		loadingRules,
		nil,
		os.Stdin,
	)
	certFile, err := kubeseal.OpenCert(
		context.Background(),
		clientConfig,
		qsealrc.ControllerNamespace,
		qsealrc.ControllerName,
		"", // the cert url
	)
	if err != nil {
		return nil, err
	}
	// #nosec: G307 -- this deferred close is fine because it is not on a writable file
	defer certFile.Close()
	pubKey, err := kubeseal.ParseKey(certFile)
	if err != nil {
		return nil, err
	}

	return &KubeSealClient{
		qsealrc:      qsealrc,
		clientConfig: clientConfig,
		pubKey:       pubKey,
	}, nil
}

// getClientSet returns a kubernetes clientset
func (k *KubeSealClient) getClientSet() (*kubernetes.Clientset, error) {
	restConfig, err := k.clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(restConfig)
}
