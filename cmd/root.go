package cmd

import (
	"fmt"
	"os"

	"github.com/davidcollom/kubectl-credentials-keychain/internal/executor"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var logger = logrus.New()

var rootCmd = &cobra.Command{
	Use:   "kubectl-credentials-keychain",
	Short: "Kubernetes credentials helper that securely stores and retrieves cluster credentials",
	Long: `kubectl-credentials-helper is a tool that helps manage Kubernetes cluster credentials
by securely storing them in the system keychain and retrieving them when needed. Credentials are never written to disk.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		r := &executor.Runner{
			Logger: logrus.StandardLogger(),
			Loader: &executor.EnvExecCredentialLoader{},
			Stdout: func(s string) { fmt.Println(s) },
		}
		return r.Run()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// Setup logrus
	logger.SetOutput(os.Stderr)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: false,
	})

	// Set debug level if environment variable is set
	if os.Getenv("KUBECTL_CREDENTIALS_KEYCHAIN_DEBUG") == "true" {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	logger.Debug("kubectl-credentials-keychain initialized")
}
