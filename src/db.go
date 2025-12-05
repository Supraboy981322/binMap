package main

import (
	"io"
	"os"
	"fmt"
	"time"
	"strings"
	"net/http"
	"github.com/charmbracelet/log"
	"github.com/Supraboy981322/gomn"
)

//might make more complex later
//  so has dedicated func
func mapDB() error {
	//shared state and mutex are needlessly
	//  complex for a simple boolean 
	for blkDB {
		time.Sleep(100 * time.Millisecond)
	};blkDB = true

	var err error
	if db, err = gomn.ReadBin(dbPath); err != nil {
		return err
	};blkDB = false

	return nil
}

//called several times
//  so has dedicated func
func updateDB(w http.ResponseWriter) error {
	//shared state and mutex are needlessly
	//  complex for a simple boolean
	for blkDB {
		time.Sleep(100 * time.Millisecond)
	};blkDB = true
	

	if !blkDB {
		if err := gomn.WrBin(db, dbPath); err != nil {
			eror(w, "failed to save db\n", err)
		};w.Write([]byte("saved db\n"))
	};blkDB = false

	return nil
}

//adds visual complexity, plus looks like spagetty,
//  so moved out of main.go into db.go as dedicated func
func dlBin(w http.ResponseWriter, typ string) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	flusher, canFlush := w.(http.Flusher)

	switch strings.ToLower(typ) {
	 case "bin", "b", "binary", "raw", "r", "gaas":
		//don't wait for it to be unlocked
		//  just set to locked and stream it
		blkDB = true

		//open db binary
		file, err := os.Open(dbPath)
		if err != nil {
			eror(w, "opening db binary", err)
		}; defer file.Close()

		//create suggested file-name header value
		t := time.Now().Format("2006-01-02_15:04:05")
		sugFileName := fmt.Sprintf("binMap_db_%s.bgomn", t)
		sugFileVal := fmt.Sprintf("attachment; filename=\"%s\"", sugFileName)

		//inform client of content
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", sugFileVal)

		// stream the binary
		if _, err = io.Copy(w, file); err != nil {
			log.Errorf("err streaming binary to client: %v", err)
		};blkDB = false

	 //this is better than the long, ugly spagetty
	 //  that it was before
	 case `key-val`, `key value`, `pair`, `text`: fallthrough
	 case `key_val`, `key-value`, `kv`, `t`, `p`: fallthrough
 	 case `key val`, `key_value`, `pairs`, `k-v`:
		//just stream basic plain-text key-value pair lines
		for key, val := range db {
			w.Write([]byte(fmt.Sprintf("%s = % x\n", key, val)))
			if canFlush { flusher.Flush() }
		}

	 case "g", "gomn", "std", "standard", "s":
		//mildly-crappy (and probably temporary)
		//  conversion to standard plain-text gomn
		for key, val := range db {
			w.Write([]byte(fmt.Sprintf("[\"%s\"] := \"% x\"\n", key, val)))
			if canFlush { flusher.Flush() }
		}

   default: //assume invalid type requested
		log.Warnf("attempt to download db as unsupported type:  %s", typ)
	}

	return
}

func deleteProd(toDefault bool) error {
	log.Warn("\033[1;4;5;31mREQUEST TO\033[0m \033[1;4;5;41mDELETE\033[0m \033[1;4;5;31mDATABASE!\033[0m")
	log.Warn("waiting \033[1;5;31m10 seconds\033[0m before deleting db")

	time.Sleep(10 * time.Second)

	log.Warn("10 SECONDS ELAPSED, STARTING deleteProd()")

	if toDefault {
		db = defDB()
	} else {
		db = gomn.Map{}
	}

	//shared state and mutex are needlessly
	//  complex for a simple boolean 
	for blkDB {
		time.Sleep(100 * time.Millisecond)
	}

	if !blkDB {
		blkDB = true
		if err := gomn.WrBin(db, dbPath); err != nil {
			log.Fatal("failed to DELETE db\n", err)
		};blkDB = false
	} 

	log.Warn("\033[1;4;5;31mDATABASE DELETED\033[0m")

	return nil
}

func clDB() {
	for true {
		time.Sleep(time.Duration(clDBSec) * time.Second)

		if blkDB { log.Debug("clDB():  db blocked, waiting until unblocked") }

		//wait for db to be unblocked
		for blkDB {
			time.Sleep(100 * time.Millisecond)
		}

		if !blkDB {		
			if clToDef { db = defDB()
			} else { db = gomn.Map{} }

			blkDB = true
			if err := gomn.WrBin(db, dbPath); err != nil {
				log.Errorf("failed to clear db:  %v", err)
			};blkDB = false
			
			log.Warn("db cleared")
		}
	}
}
