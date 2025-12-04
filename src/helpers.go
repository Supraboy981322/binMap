package main

import (
	"os"
	"fmt"
	"net/http"
	"github.com/charmbracelet/log"
)

func getKey(r *http.Request) string {
	return chkHeaders(r, []string{"k", "key"})
}

func getVal(r *http.Request) string {
	return chkHeaders(r, []string{"v", "val", "value"})
}

func chkHeaders(r *http.Request, list []string) string {
	var res string
	for _, chk := range list {
		res = r.Header.Get(chk)
		if res != "" { break }
	}
	return res
}

func eror(w http.ResponseWriter, str string, err error) {
	erorr := fmt.Sprintf("%s:  %v", str, err)
	log.Error(erorr)
	w.Write([]byte(erorr+"\n"))
}

//looks like spagetty because of ansi color codes
func logReq(p string, ip string, extra string) {
	if canLog, ok := config["log requests"].(bool); ok && canLog {
		log.Printf("\033[1;36m[req]\033[0m "+
			"page=\033[1;37m%s\033[0m ; "+
			"ip=\033[1;37m%s\033[0m ; %s", p, ip, extra)
	}
}

func defConf() {
	log.Warnf("creating default config, "+
			"\033[1;37mpath=\033[0m\033[1;32m%s\033[0m", configPath)

	//write the file
	err := os.WriteFile(configPath, defConfig(), 0666) //rw permission
	if err != nil {
		log.Fatal("creating default config:  %v", err)
	}

	//quit with msg to read config
	log.Fatalf("please review default config "+
			"(at \033[1;4;5;32m%s\033[0m), "+
			"and restart the server", configPath)
}
