package cmd

import (
	"github.com/davidcollom/kubectl-credentials-keychain/internal/undo"
	"github.com/spf13/cobra"
)

var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "This makes your kubeconfig insecure!",
	RunE: func(cmd *cobra.Command, args []string) error {
		kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
		user, _ := cmd.Flags().GetString("user")

		r := &undo.Runner{
			KubeconfigPath: kubeconfig,
			SpecificUser:   user,
			FS:             OsFs,
			Logger:         logger,
		}

		return r.Run()
	},
}

func init() {
	undoCmd.Flags().StringP("kubeconfig", "c", "", "Kubeconfig path")
	undoCmd.Flags().StringP("user", "u", "", "Undo for specific user only")

	rootCmd.AddCommand(undoCmd)
}
