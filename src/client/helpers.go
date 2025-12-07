package main

import (
	"os"
	"fmt"
	"slices"
	"errors"
	"net/url"
	"path/filepath"
	"github.com/Supraboy981322/gomn"
)

func parseArgs() {
	var taken []int
	args := os.Args[1:]
	for i, arg := range args {
		if !slices.Contains(taken, i) {
			if arg[0] == '-' && arg[1] != '-' {
				for _, a := range arg[1:] {
					switch a {
           case 'i':
						if ok := chkAhead(args, i); ok {
							input = args[i+1]
							if input == "-" { useStdin = true }
							taken = append(taken, i+1)
						} else {
							err := errors.New("input arg requires a value")
							eror("no input provided", err)
						}
					 case 'k':
						if ok := chkAhead(args, i); ok {
							key = args[i+1]
							taken = append(taken, i+1)
						} else {
							err := errors.New("input arg requires a value")
							eror("no value provided", err)
						}
					 case 'B': binary = true
					 case 'v':
						if ok := chkAhead(args, i); ok {
							val = []byte(args[i+1])
							taken = append(taken, i+1)
						} else {
							err := errors.New("value arg requires a value")
							eror("no value provided", err)
						}
           case 'S':
						if ok := chkAhead(args, i); ok {
						  addr = args[i+1]
							taken = append(taken, i+1)
						} else {
							err := errors.New("server arg requires a value")
							eror("no value provided", err)
						}
					 case 's':
						if act == "" { act = "set"
						} else {
							err := fmt.Errorf("action set to %s", act)
							eror("only one action allowed", err) }
					 case 'g':
						if act == "" { act = "get"
						} else {
							err := fmt.Errorf("action set to %s", act)
							eror("only one action allowed", err) }
					 case 'D':
						if act == "" { act = "del"
						} else {
							err := fmt.Errorf("action set to %s", act)
							eror("only one action allowed", err) }
					 case 'V': verbose = true
					 default:
					  eror("invalid arg", errors.New(string(a)))
					}
				}
			} else {
				switch arg[2:] {
 				 case "server":
					if ok := chkAhead(args, i); ok {
					  addr = args[i+1]
						taken = append(taken, i+1)
					} else {
						err := errors.New("server arg requires a value")
						eror("no value provided", err)
					}
				 case "key":
					if ok := chkAhead(args, i); ok {
						key = args[i+1]
						taken = append(taken, i+1)
					} else {
						err := errors.New("key arg requires a value")
						eror("no value provided", err)
					}
				 case "val", "value":
					if ok := chkAhead(args, i); ok {
						val = []byte(args[i+1])
						taken = append(taken, i+1)
					} else {
						err := errors.New("value arg requires a value")
						eror("no value provided", err)
					}
		     default:
				  eror("invalid arg", errors.New(arg))
				}
			}
		}
	}
}

func mkDefConf() {
	file := []byte(`//binMap client config
["server address"] := "http://[::1]:4780"
["verbose"] := false`)

	if err := os.MkdirAll(filepath.Dir(confPath), 0644); err != nil {
		eror("failed to make config path", err)
	}
	
	if err := os.WriteFile(confPath, file, 0644); err != nil {
		eror("failed to write default config", err)
	}
}

func chkAhead(arr []string, i int) bool {
	if len(arr) > i+1 {
		return true
	}
	return false
}

func parseConfig() {
	var ok bool
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

	if addr == "" {
  	if addr, ok = conf["server address"].(string); !ok {
			err = errors.New("not a string")
			eror("failed to assert server address", err)
		} else { verbLog("asserted server address") }
	} else { verbLog("address already set, not checking config") }

	if addr, ok = validateURL(addr); !ok {
		err := fmt.Errorf("invalid url: '%s'", addr)
		eror("failed to parse url", err)
	} else { verbLog("validated url") }

	if !verbose {
		if verbose, ok = conf["verbose"].(bool); !ok {
			err = errors.New("not a bool")
			eror("failed to assert \"verbose\" in config", err)
		} else { verbLog("asserted verbosity") }
	} else { verbLog("verbosity already set") }

	if string(val) == "-" { useStdin = true }

	if act == "" { act = "get" }

	var hasIn bool
	itms := []bool{
		val != nil,
		input != "",
	};for _, chk := range itms {
		if chk && hasIn	{
			err := errors.New("only accepts one value")
			eror("too many values", err)
		} else if chk { hasIn = true }
	}
}

func validateURL(og string) (string, bool) {
	u, err := url.ParseRequestURI(og)
	if err != nil { return og, false }

	if u.Scheme == "" { u.Scheme = "https" }

	return u.String(), true
}
