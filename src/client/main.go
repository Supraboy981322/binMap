package main

import (
	"os"
	"io"
	"fmt"
	"time"
	"bytes"
	"errors"
	"net/http"
	"io/ioutil"
	"github.com/Supraboy981322/gomn"
)

var (
	act string
	val []byte
	key string
	timeout int
	addr string
	input string
	verbose bool
	useStdin bool
	conf gomn.Map
	confPath string
)

func init() {
	parseArgs()
	parseConfig()
}

func main() {
	if verbose {
		foo := map[string]string {
			"server": addr,
			"verbose": fmt.Sprint(verbose),
			"useStdin": fmt.Sprint(useStdin),
			"confPath": confPath,
			"key": key,
			"val": string(val),
			"input": input,
		};for k, v := range foo {
			verbLog(fmt.Sprintf("%s: %s", k, v))
		}
	}
	mkReq()
}

func eror(str string, err error) {
	fmt.Fprintf(os.Stderr, "\033[1;31merr %s:\033[0m\n"+
			"\t\033[1;41;30m%v\033[0m\n", str, err)
	os.Exit(1)
}

func verbLog(str string) {
	if verbose {
		fmt.Printf("\033[1;36m[V]:\033[0m  %s\n", str)
	}
}

func mkReq() {
	if addr == "" {
		err := errors.New("address is empty")
		eror("can't send request to server", err)
	}

	timeDir := time.Duration(timeout)
	client := &http.Client{
		Timeout: time.Second * timeDir,
	};_ = client

	var payload io.Reader
	if act == "set" { 
		var err error
		if input != "" {
			payload, err = os.Open(input)
			if err != nil { eror("failed to read file", err) }
		} else if val != nil {
			payload = bytes.NewReader(val)
		} else { eror("FATAL: UNCAUGHT ERR", errors.New("value")) }
	}

	url := fmt.Sprintf("%s/%s", addr, act)
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		eror("failed to create request", err)
	}

	req.Header.Add("key", key)

	resp, err := client.Do(req)
	if err != nil {
		eror("failed to send request", err)
	};defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		eror("failed to read response body", err)
	};defer resp.Body.Close()

	print(string(body))
}
