package main

import (
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

func logReq(p string, ip string, ext string) {
	log.Printf("\033[1;36m[req]\033[0m "+
		"page=\033[1;37m%s\033[0m ; "+
		"ip=\033[1;37m%s\033[0m ; %s", p, ip, ext)
}
