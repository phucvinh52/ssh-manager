package sshconfig

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type SSHHost struct {
	HostName     string
	User         string
	Port         string
	IdentityFile string
}
type SSHHostFull struct {
	Host string
	*SSHHost
}
type SSHConfig struct {
	Hosts map[string]*SSHHostFull
}

func (s *SSHConfig) Filter(txSearch string) []*SSHHostFull {
	var arr []*SSHHostFull
	txSearch = strings.ToLower(txSearch)
	for _, v := range s.Hosts {
		hLower := strings.ToLower(v.Host)
		hnLower := strings.ToLower(v.HostName)

		if !strings.Contains(hLower, txSearch) && !strings.Contains(hnLower, txSearch) {
			continue
		}

		arr = append(arr, v)
	}
	return arr
}

type OptParseSSHConfig struct {
	Path string
}

func ParseSSHConfig(opts *OptParseSSHConfig) (*SSHConfig, error) {
	if opts == nil {
		opts = new(OptParseSSHConfig)
	}
	if len(opts.Path) == 0 {
		opts.Path = filepath.Join(os.Getenv("HOME"), ".ssh", "config")
	}

	osFile, err := os.Open(opts.Path)
	if err != nil {
		return nil, fmt.Errorf("ParseSSHConfig: %w", err)
	}
	readerFile := bufio.NewReader(osFile)

	var sshCfg = new(SSHConfig)
	sshCfg.Hosts = make(map[string]*SSHHostFull)
	var lastSSHHost = new(SSHHost)

	for {
		lineB, _, err := readerFile.ReadLine()
		if err != nil {
			break
		}
		lineS := strings.TrimSpace(string(lineB))
		parts := strings.Split(lineS, " ")

		switch parts[0] {
		case "Host":
			lastSSHHost = new(SSHHost)
			for _, v := range parts {
				sshCfg.Hosts[v] = new(SSHHostFull)
				sshCfg.Hosts[v].Host = v
				sshCfg.Hosts[v].SSHHost = lastSSHHost
			}
		case "HostName":
			lastSSHHost.HostName = parts[1]
		case "Port":
			lastSSHHost.Port = parts[1]
		case "User":
			lastSSHHost.User = parts[1]
		case "IdentityFile":
			lastSSHHost.IdentityFile = parts[1]
		}
	}
	return sshCfg, nil
}
