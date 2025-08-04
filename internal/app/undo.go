package app

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/davidcollom/kubectl-credentials-keychain/internal/keychain"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	clientauthentication "k8s.io/client-go/pkg/apis/clientauthentication/v1"
	"k8s.io/client-go/tools/clientcmd"
	cmdApi "k8s.io/client-go/tools/clientcmd/api"
)

type Runner struct {
	KubeconfigPath string
	SpecificUser   string
	FS             FileSystem
	Keychain       keychain.Keychain
	Prompter       Prompter
	Logger         *logrus.Logger
}

func (r *Runner) Secure() error {
	path, err := r.resolveKubeConfig()
	if err != nil {
		return err
	}

	cfg, err := r.loadAndBackup(path)
	if err != nil {
		return err
	}

	executable, _ := os.Executable()

	for name, user := range cfg.AuthInfos {
		if r.SpecificUser != "" && r.SpecificUser != name {
			continue
		}
		if !hasSensitive(user) {
			continue
		}
		r.Logger.Infof("Securing user: %s", name)

		minimalCfg, err := extractUserConfig(cfg, name)
		if err != nil {
			return err
		}
		b64Cfg := base64.StdEncoding.EncodeToString(minimalCfg)

		for ctxName, ctx := range cfg.Contexts {
			log := r.Logger.WithField("context", ctxName)
			if ctx.AuthInfo != name {
				continue
			}
			cluster := cfg.Clusters[ctx.Cluster]
			ok, err := r.Prompter.Confirm(fmt.Sprintf("Create secret for %s (%s)?", ctx.Cluster, cluster.Server))
			if err != nil || !ok {
				continue
			}

			err = r.Keychain.CreateSecret(ctx.Cluster, cluster.Server, b64Cfg)
			if err == keychain.ErrorDuplicateItem {
				replace, err := r.Prompter.Confirm(fmt.Sprintf("Secret %s already exists. Replace it?", ctx.Cluster))
				if err != nil || !replace {
					continue
				}
				if err := r.Keychain.DeleteSecret(cluster.Server); err != nil {
					return err
				}
				if err := r.Keychain.CreateSecret(ctx.Cluster, cluster.Server, b64Cfg); err != nil {
					return err
				}
			} else if err != nil {
				return err
			}
			log.Infof("Secret created: %s", ctx.Cluster)
		}

		if ok, err := r.Prompter.Confirm(fmt.Sprintf("Remove sensitive parts from user %s?", name)); err == nil && ok {
			sanitizeUser(user, executable)
		}
	}

	return clientcmd.WriteToFile(cfg, path)
}

// --- helpers ---

func (r *Runner) resolveKubeConfig() (string, error) {
	if r.KubeconfigPath != "" {
		return filepath.Abs(r.KubeconfigPath)
	}
	if val, ok := os.LookupEnv("KUBECONFIG"); ok {
		return filepath.Abs(val)
	}
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Abs(filepath.Join(home, ".kube", "config"))
}

func (r *Runner) loadAndBackup(path string) (cmdApi.Config, error) {
	cfg, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: path},
		&clientcmd.ConfigOverrides{},
	).RawConfig()
	if err != nil {
		return cmdApi.Config{}, err
	}
	_ = clientcmd.WriteToFile(cfg, path+".back")
	return cfg, nil
}

func extractUserConfig(cfg cmdApi.Config, name string) ([]byte, error) {
	auth := cfg.AuthInfos[name]
	tmpCfg := cmdApi.Config{
		AuthInfos: map[string]*cmdApi.AuthInfo{name: {
			ClientCertificateData: auth.ClientCertificateData,
			ClientKeyData:         auth.ClientKeyData,
			Username:              auth.Username,
			Password:              auth.Password,
		}},
	}
	return clientcmd.Write(tmpCfg)
}

func hasSensitive(user *cmdApi.AuthInfo) bool {
	return len(user.ClientCertificateData) > 0 || len(user.ClientKeyData) > 0 || user.Username != "" || user.Password != ""
}

func sanitizeUser(user *cmdApi.AuthInfo, execPath string) {
	user.ClientCertificateData = nil
	user.ClientKeyData = nil
	user.Username = ""
	user.Password = ""
	user.Exec = &cmdApi.ExecConfig{
		APIVersion:         clientauthentication.SchemeGroupVersion.String(),
		Command:            execPath,
		ProvideClusterInfo: true,
		InteractiveMode:    cmdApi.NeverExecInteractiveMode,
	}
}
