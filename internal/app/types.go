package app

import (
	"errors"

	"github.com/charmbracelet/huh"
	"github.com/davidcollom/kubectl-credentials-keychain/internal/keychain"
	"github.com/sirupsen/logrus"
	clientauthentication "k8s.io/client-go/pkg/apis/clientauthentication/v1"
	"k8s.io/client-go/tools/auth/exec"
)

type FileSystem interface {
	Exists(path string) (bool, error)
	IsDir(path string) (bool, error)
}

type Prompter interface {
	Confirm(title string) (bool, error)
}

type HuhPrompter struct{}

func (p *HuhPrompter) Confirm(title string) (bool, error) {
	var result bool
	err := huh.NewConfirm().Title(title).Value(&result).Run()
	return result, err
}

type ExecCredentialLoader interface {
	Load() (*clientauthentication.ExecCredential, error)
}

type EnvExecCredentialLoader struct{}

func (l *EnvExecCredentialLoader) Load() (*clientauthentication.ExecCredential, error) {
	ec, _, err := exec.LoadExecCredentialFromEnv()
	if err != nil {
		return nil, err
	}
	cast, ok := ec.(*clientauthentication.ExecCredential)
	if !ok {
		return nil, errors.New("failed to cast ExecCredential")
	}
	return cast, nil
}

type Runner struct {
	KubeconfigPath string
	SpecificUser   string
	FS             FileSystem
	Keychain       keychain.Keychain
	Prompter       Prompter
	Logger         *logrus.Logger
}
