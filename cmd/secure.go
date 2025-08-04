package cmd

import (
	"github.com/davidcollom/kubectl-credentials-keychain/internal/secure"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// OsFs is an instance of afero.NewOsFs
var OsFs = afero.Afero{Fs: afero.NewOsFs()}

// FileExists afero for some reason does not have such a function, so...
func FileExists(filename string) (bool, error) {
	e, err := OsFs.Exists(filename)
	if err != nil {
		return e, err
	}

	e, err = OsFs.IsDir(filename)
	if err != nil {
		return e, err
	}

	return !e, nil
}

func init() {
	secureCmd.Flags().StringP("kubeconfig", "c", "", "Kubeconfig path")
	secureCmd.Flags().StringP("user", "u", "", "Secure specific user instead of all")

	rootCmd.AddCommand(secureCmd)
}

var secureCmd = &cobra.Command{
	Use:   "secure",
	Short: "This makes your kubeconfig secure!",
	RunE: func(cmd *cobra.Command, args []string) error {
		kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
		user, _ := cmd.Flags().GetString("user")

		r := &secure.Runner{
			KubeconfigPath: kubeconfig,
			SpecificUser:   user,
			FS:             OsFs,
			Prompter:       &secure.HuhPrompter{},
			Logger:         logger,
		}

		return r.Run()
	},
}
