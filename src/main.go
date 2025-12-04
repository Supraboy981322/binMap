package main

import (
	"io"
	"errors"
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

	if err = initDB(); err != nil {
		log.Fatalf("failed to initialize database:  %v", err)
	}

	if err = mapDB(); err != nil {
		log.Fatalf("failed to map database:  %v", err)
	}
}

func main() {
	log.Info("started")
	http.HandleFunc("/get", getHan)
	http.HandleFunc("/set", setHan)
	http.HandleFunc("/db", dbHan)
	
	portStr := ":"+strconv.Itoa(port)
	log.Infof("listening on port %d", port)
	log.Fatal(http.ListenAndServe(portStr, nil))
}

func getHan(w http.ResponseWriter, r *http.Request) {
	var ok bool

	log.Infof("[req]:  /get  %s", r.RemoteAddr)

	var key string
	if key = getKey(r); key == "" {
		bod, err := io.ReadAll(r.Body)
		if err != nil {
			eror(w, "reading req body", err)
			return
		};key = string(bod)
		if key == "" {
			w.Write([]byte("need key\n"))
			return
		}
	}

	var val []byte
	if val, ok = db[key].([]byte); !ok {
		if db[key] == nil {
			err := errors.New(key)
			eror(w, "key does not exist", err)
		}	else {
			err := errors.New(key)
			eror(w, "invalid value in db. key", err)
		}
		return
	}

	newline := chkHeaders(r, []string{"n", "newline", `\n`, "\n"})
	if newline == "" || newline == "true" {
		val = append(val, []byte("\n")...)
	}

	w.Write(val)
}

func setHan(w http.ResponseWriter, r *http.Request) {
	log.Info("[req]:  /set  %s", r.RemoteAddr)
	
	var key string
	if key = getKey(r); key == "" {
		if key == "" {
			w.Write([]byte("need key\n"))
			return
		}
	}
	
	var val string
	if val = getVal(r); val == "" {
		bod, err := io.ReadAll(r.Body)
		if err != nil {
			eror(w, "reading req body", err)
			return
		};val = string(bod)
		if val == "" {
			w.Write([]byte("need value\n"))
			return
		}
	}

	db[key] = []byte(val)
	w.Write([]byte("added to db\n"))

	
	if err := gomn.WrBin(db, dbPath); err != nil {
		eror(w, "failed to save db\n", err)
	};w.Write([]byte("saved db\n"))

	w.Write([]byte("done\n"))
}

func dbHan(w http.ResponseWriter, r *http.Request) {
	typ := chkHeaders(r, []string{"type", "t", "typ"})
	if typ == "" {
		typ = "bin"
	}

	dlBin(w, typ)
}
