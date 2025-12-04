package main

import (
	"io"
	"fmt"
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
	http.HandleFunc("/del", delHan)
	
	portStr := ":"+strconv.Itoa(port)
	log.Infof("listening on port %d", port)
	log.Fatal(http.ListenAndServe(portStr, nil))
}

func getHan(w http.ResponseWriter, r *http.Request) {
	var ok bool

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
	
	logReq("/get", r.RemoteAddr, "key="+key)

	w.Write(val)
}

func setHan(w http.ResponseWriter, r *http.Request) {
	
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
	
	logReq("/set", r.RemoteAddr, "key="+key)

	db[key] = []byte(val)
	w.Write([]byte("added to db\n"))

	updateDB(w)

	w.Write([]byte("done\n"))
}

func delHan(w http.ResponseWriter, r *http.Request) {
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

	logReq("/del", r.RemoteAddr, "key="+key)

	delete(db, key)
	updateDB(w)
	
	w.Write([]byte(fmt.Sprintf("deleted:  %s\n", key)))
}

func dbHan(w http.ResponseWriter, r *http.Request) {
	typ := chkHeaders(r, []string{"type", "t", "typ"})
	if typ == "" {
		bod, err := io.ReadAll(r.Body)
		if err != nil {
			eror(w, "reading req body", err)
			return
		};typ = string(bod)
		if typ == "" {
			typ = "bin"
		}
	}

	logReq("/db", r.RemoteAddr, "typ="+typ)
	
	dlBin(w, typ)
}
