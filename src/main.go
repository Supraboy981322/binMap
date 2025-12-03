package main

import (
	"fmt"
	"strconv"
	"net/http"
	"github.com/charmbracelet/log"
	"github.com/Supraboy981322/gomn"
)

var (
	db gomn.Map
	port = 8944
	dbPath string
	logLvl string
	config gomn.Map
	configPath string
	bin map[string][]byte
)

func init() {
	var err error

	log.SetLevel(log.DebugLevel)

	configPath = "conf.gomn"
	if err = configure(); err != nil {
		log.Fatalf("failed to configure:  %v", err)
	}

	if err = mapDB(); err != nil {
		log.Fatalf("failed to map database:  %v", err)
	}
}

func main() {
	log.Info("started")
	http.HandleFunc("/get", getHandler)
	
	portStr := ":"+strconv.Itoa(port)
	log.Infof("listening on port %d", port)
	log.Fatal(http.ListenAndServe(portStr, nil))
}

func getHan(w http.ResponseWriter, r *http.Request) {
	var ok bool

	log.Infof("[req]:  /get  %s", r.RemoteAddr)
	var key string
	for _, chk := range []string{"k", "key"} {
		key = r.Header.Get(chk)
		if key != "" { break }
	};if key == "" {
		w.Write([]byte("need key\n"))
		return
	}

	var val []byte
	if val, ok = db[key].([]byte); !ok {
		eror := fmt.Sprintf("invalid value in db. key:  %s", key)
		log.Error(eror)
		w.Write([]byte(eror))
		return
	}

	val = append(val, []byte("\n")...)
	w.Write(val)
}

