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
	bin map[string][]byte
	configPath = "conf.gomn"
)

func init() {
	var err error

	log.Info("initializing server...")

	//temporarilly use debug mode
	//  (changed when parsing config)
	log.SetLevel(log.DebugLevel)

	//configure the server (clearly)
	if err = configure(); err != nil {
		log.Fatalf("failed to configure:  %v", err)
	}

	//initialize database (clearly)
	if err = initDB(); err != nil {
		log.Fatalf("failed to initialize database:  %v", err)
	}

	//generate database map in memory
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
	
	log.Infof("listening on port %d", port)

	//looks cleaner 
	portStr := ":"+strconv.Itoa(port)
	log.Fatal(http.ListenAndServe(portStr, nil))
}

func getHan(w http.ResponseWriter, r *http.Request) {
	var ok bool //avoids golang quirk (is bug?)

	//get the key
	var key string //used later, so init outside block
	if key = getKey(r); key == "" {
		//read the req body if no key header
		bod, err := io.ReadAll(r.Body)
		if err != nil {
			eror(w, "reading req body", err)
			return
		};key = string(bod)
		//return err if still no key 
		if key == "" {
			w.Write([]byte("need key\n"))
			return
		}
	}

	//get val from db 
	var val []byte //used later, so init outside block
	if val, ok = db[key].([]byte); !ok {
		//check problem
		if db[key] != nil {//if not []byte, it's invalid 
			err := errors.New(key)
			eror(w, "invalid value in db. key", err)
		}	else {//if empty, assume missing
			err := errors.New(key)
			eror(w, "key does not exist", err)
		}
		return
	}

	//check if header specifies to not use newline
	newline := chkHeaders(r, []string{"n", "newline", `\n`, "\n"})
	if newline == "" || newline == "true" {
		val = append(val, []byte("\n")...)
	}
	
	//self-explainitory
	logReq("/get", r.RemoteAddr, "key="+key)

	//send value
	w.Write(val)
}

func setHan(w http.ResponseWriter, r *http.Request) {
	//get the key header
	var key string//used later, so init outside block
	if key = getKey(r); key == "" {
		if key == "" {
			w.Write([]byte("need key\n"))
			return
		}
	}
	
	//get the value header
	var val string//used later, so init outside block
	if val = getVal(r); val == "" {
		//if no val header, use body
		bod, err := io.ReadAll(r.Body)
		if err != nil {
			eror(w, "reading req body", err)
			return
		};val = string(bod)
		//if body is also empty, return err
		if val == "" {
			w.Write([]byte("need value\n"))
			return
		}
	}
	
	//self-explainitory
	logReq("/set", r.RemoteAddr, "key="+key)

	//set value in db (in memory)
	db[key] = []byte(val)
	w.Write([]byte("added to db\n"))

	updateDB(w) //save in-memory db to disk

	//confirm completion
	//  (updateDB sends progress to client)
	w.Write([]byte("done\n"))
}

func delHan(w http.ResponseWriter, r *http.Request) {
	//get the key
	var key string//used later, so init outside block
	if key = getKey(r); key == "" {
		//use body if header empty
		bod, err := io.ReadAll(r.Body)
		if err != nil {
			eror(w, "reading req body", err)
			return
		};key = string(bod)
		//if no body, err
		if key == "" {
			w.Write([]byte("need key\n"))
			return
		}
	}

	//self-explainitory
	logReq("/del", r.RemoteAddr, "key="+key)

	delete(db, key) //delete from in-memory db
	updateDB(w) //save changes to disk

	//confirm completion
	msg := fmt.Sprintf("deleted:  %s\n", key)
	w.Write([]byte(msg))
}

func dbHan(w http.ResponseWriter, r *http.Request) {
	//get the type header 
	typ := chkHeaders(r, []string{"type", "t", "typ"})
	if typ == "" {
		//check body
		bod, err := io.ReadAll(r.Body)
		if err != nil {
			eror(w, "reading req body", err)
			return
		};typ = string(bod)
		//if no body, err
		if typ == "" {
			typ = "bin"
		}
	}

	//self-explainitory
	logReq("/db", r.RemoteAddr, "typ="+typ)
	
	//stream db to client
	dlBin(w, typ)
}
