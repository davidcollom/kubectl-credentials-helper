package cmd

import (
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func init() {
	rootCmd.AddCommand(docsCmd)
}

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generate docs",
	Long:  `Uses Cobra to generate CLI docs`,
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		// Check to ensure docs dir exists, if not (attempt to) create one
		docsDir := path.Join(cwd, "docs")
		if exists, err := OsFs.DirExists(docsDir); err != nil {
			log.Fatal(err)
		} else if !exists {
			err = OsFs.MkdirAll(docsDir, 0755)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = doc.GenMarkdownTree(rootCmd, path.Join(cwd, "docs"))
		if err != nil {
			log.Fatal(err)
		}
	},
}
