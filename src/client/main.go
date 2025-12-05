package main

import (
	"os"
	"fmt"
	"errors"
//	"net/http"
	"path/filepath"
	"github.com/Supraboy981322/gomn"
)

var (
	val []byte
	key string
	addr string
	input string
	verbose bool
	useStdin bool
	conf gomn.Map
	confPath string
)

func init() {
	var ok bool
	parseArgs()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		eror("failed to get user home dir (for config)", err)
	} else { verbLog("got home dir") }

	confPath = filepath.Join(homeDir,
			".config/Supraboy981322/config.gomn")

	if conf, err = gomn.ParseFile(confPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("no config, creating default")
			mkDefConf()
		} else { eror("failed to parse config", err) }
	} else { verbLog("parsed config") }

	if addr != "" {
  	if addr, ok = conf["server address"].(string); !ok {
			err = errors.New("not a string")
			eror("failed to assert server address", err)
		} else { verbLog("asserted server address") }
	} else { verbLog("address already set, not checking config") }
}

func main() {
	foo := map[string]string {
		"server": addr,
		"verbose": fmt.Sprint(verbose),
		"useStdin": fmt.Sprint(useStdin),
		"confPath": confPath,
		"key": key,
		"val": string(val),
		"input": input,
	}
	for k, v := range foo {
		fmt.Printf("%s: %s\n", k, v)
	}
}

func eror(str string, err error) {
	fmt.Fprintf(os.Stderr, "%s:  %v\n", str, err)
	os.Exit(1)
}

func verbLog(str string) {
	if verbose { fmt.Println(str) }
}
