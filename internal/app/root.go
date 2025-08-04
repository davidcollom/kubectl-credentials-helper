package app

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	clientauthentication "k8s.io/client-go/pkg/apis/clientauthentication/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func (r *Runner) RunRoot() error {
	r.Logger.Debugf("KUBERNETES_EXEC_INFO: %q", os.Getenv("KUBERNETES_EXEC_INFO"))

	ec, err := r.Loader.Load()
	if err != nil {
		return fmt.Errorf("failed to load exec credential: %w", err)
	}

	clusterEndpoint := ec.Spec.Cluster.Server
	if clusterEndpoint == "" {
		return errors.New("empty cluster endpoint")
	}
	r.Logger.Debugf("Cluster endpoint: %s", clusterEndpoint)

	secretName, secretB64, err := r.Keychain.GetSecret(clusterEndpoint)
	if err != nil {
		return fmt.Errorf("failed to get secret: %w", err)
	}
	r.Logger.Debugf("Found secret: %s", secretName)

	decoded, err := base64.StdEncoding.DecodeString(secretB64)
	if err != nil {
		return fmt.Errorf("failed to decode base64: %w", err)
	}

	cfg, err := clientcmd.Load(decoded)
	if err != nil {
		return fmt.Errorf("failed to parse kubeconfig: %w", err)
	}

	authInfo := cfg.AuthInfos[secretName]
	if authInfo == nil {
		return fmt.Errorf("auth info for %s not found", secretName)
	}

	ec.APIVersion = clientauthentication.SchemeGroupVersion.String()
	ec.Status = &clientauthentication.ExecCredentialStatus{
		ClientCertificateData: string(authInfo.ClientCertificateData),
		ClientKeyData:         string(authInfo.ClientKeyData),
	}

	out, err := json.Marshal(ec)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	r.Stdout(string(out))
	return nil
}
