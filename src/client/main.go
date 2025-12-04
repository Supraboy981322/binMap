package main

import (
	"os"
	"net/http"
	"path/filepath"
	"github.com/Supraboy981322/gomn"
)

var (
	conf gomn.Map
	confPath string
)

func init() {
	homeDir, err := os.UserHomeDir()
	confPath := filepath.Join(homeDir, ".config/Supraboy981322/config.gomn")

	if conf, err = gomn.ParseFile(confPath); err != nil {
		eror("failed to parse config", err)
	}
}

func eror(str string, err error)
	fmt.Fprintf(os.Stderr, "%s:  %v", str, err)
	os.Exit(1)
}
