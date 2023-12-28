package main

import (
	"github.com/phucvinh52/ssh-manager/internal/appui"
	"github.com/phucvinh52/ssh-manager/pkg/sshconfig"
)

func main() {
	sshCfg, err := sshconfig.ParseSSHConfig(nil)
	if err != nil {
		panic(err)
	}
	appui.CreateApp(sshCfg).Start()
}
