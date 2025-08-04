package app

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"k8s.io/client-go/tools/clientcmd"
	cmdApi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/davidcollom/kubectl-credentials-keychain/internal/keychain"
)

type FileSystem interface {
	Exists(string) (bool, error)
	IsDir(string) (bool, error)
}

type Keychain interface {
	GetSecret(service string) (name string, value string, err error)
}

type Logger interface {
	Infof(string, ...interface{})
	Debugf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
}

type Runner struct {
	KubeconfigPath string
	SpecificUser   string
	FS             FileSystem
	Keychain       Keychain
	Logger         Logger
}

func (r *Runner) RunSecure() error {
	path, err := r.resolveKubeConfig()
	if err != nil {
		return err
	}

	cfg, err := loadAndBackup(path)
	if err != nil {
		return err
	}

	executable, _ := os.Executable()

	for name, user := range cfg.AuthInfos {
		if r.SpecificUser != "" && r.SpecificUser != name {
			r.Logger.Debugf("Skipping user %s: not target user %s", name, r.SpecificUser)
			continue
		}

		if user.Exec == nil || !strings.HasSuffix(user.Exec.Command, executable) {
			r.Logger.Debugf("Skipping user %s: exec not matching this tool", name)
			continue
		}

		r.Logger.Infof("Restoring credentials for user: %s", name)

		for _, ctx := range cfg.Contexts {
			if ctx.AuthInfo != name {
				continue
			}

			cluster := cfg.Clusters[ctx.Cluster]
			if cluster == nil {
				continue
			}

			secretName, secretB64, err := r.Keychain.GetSecret(cluster.Server)
			if err != nil {
				if err == keychain.ErrorItemNotFound {
					continue
				}
				return err
			}

			decoded, err := base64.StdEncoding.DecodeString(secretB64)
			if err != nil {
				return err
			}

			secretCfg, err := clientcmd.Load(decoded)
			if err != nil {
				return err
			}

			restored := secretCfg.AuthInfos[secretName]
			user.Username = restored.Username
			user.Password = restored.Password
			user.ClientCertificateData = restored.ClientCertificateData
			user.ClientKeyData = restored.ClientKeyData
			user.Exec = nil

			r.Logger.Infof("Restored from secret: %s (%s)", ctx.Cluster, cluster.Server)
			break
		}
	}

	err = clientcmd.WriteToFile(cfg, path)
	if err != nil {
		return err
	}

	r.Logger.Infof("Wrote updated kubeconfig to: %s", path)
	return nil
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
	path := filepath.Join(home, ".kube", "config")
	return filepath.Abs(path)
}

func loadAndBackup(path string) (cmdApi.Config, error) {
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
