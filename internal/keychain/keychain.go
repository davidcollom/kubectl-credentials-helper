package keychain

const (
	Service     = "Kubernetes"
	AccessGroup = "github.com/plumber-cd/kubectl-credentials-helper"
)

type Keychain interface {
	CreateSecret(name, server, data string) error
	DeleteSecret(server string) error
	GetSecret(server string) (name string, value string, err error)
}
