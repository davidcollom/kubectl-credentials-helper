package main

import (
	"log"

	"github.com/davidcollom/kubectl-credentials-keychain/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
