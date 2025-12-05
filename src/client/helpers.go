package main

import (
	"os"
	"slices"
	"errors"
	"path/filepath"
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
					 case 'v':
						if ok := chkAhead(args, i); ok {
							val = []byte(args[i+1])
							taken = append(taken, i+1)
						} else {
							err := errors.New("value arg requires a value")
							eror("no value provided", err)
						}
           case 's':
						if ok := chkAhead(args, i); ok {
						  addr = args[i+1]
							taken = append(taken, i+1)
						} else {
							err := errors.New("server arg requires a value")
							eror("no value provided", err)
						}
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
