package qseal

import (
	"context"
	"crypto/rsa"
	"fmt"

	"github.com/bitnami-labs/sealed-secrets/pkg/crypto"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/keyutil"
)

const labelSelector = "sealedsecrets.bitnami.com/sealed-secrets-key"

// RetrieveKeys retrieves the keys from the Kubernetes cluster
// the map is fingerprint -> private key
func (k *KubeSealClient) RetrievePrivateKeys() (map[string]*rsa.PrivateKey, error) {
	clientset, err := k.getClientSet()
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %v", err)
	}

	secrets, err := clientset.CoreV1().
		Secrets(k.qsealrc.ControllerNamespace).
		List(context.TODO(), metav1.ListOptions{
			LabelSelector: labelSelector,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}
	// map[fingerprint]key
	keys := make(map[string]*rsa.PrivateKey)
	for _, secret := range secrets.Items {
		tlsKey, ok := secret.Data["tls.key"]
		if !ok {
			return nil, fmt.Errorf("secret must contain a 'tls.data' key")
		}
		pkey, err := parsePrivKey(tlsKey)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		fingerprint, err := crypto.PublicKeyFingerprint(&pkey.PublicKey)
		if err != nil {
			return nil, fmt.Errorf("failed to get fingerprint: %w", err)
		}
		keys[fingerprint] = pkey
	}
	return keys, nil
}

func parsePrivKey(b []byte) (*rsa.PrivateKey, error) {
	key, err := keyutil.ParsePrivateKeyPEM(b)
	if err != nil {
		return nil, err
	}
	switch rsaKey := key.(type) {
	case *rsa.PrivateKey:
		return rsaKey, nil
	default:
		return nil, fmt.Errorf("unexpected private key type %T", key)
	}
}