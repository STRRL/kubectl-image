package main

import (
	"github.com/STRRL/kubectl-push/pkg/cmd"
)

func main() {
	cmd.NewCmdPush().Execute()

}
