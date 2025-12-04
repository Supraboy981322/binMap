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

func mapDB() error {
	var err error
	if db, err = gomn.ReadBin(dbPath); err != nil {
		return err
	}

	return nil
}

func updateDB(w http.ResponseWriter) error {
	if err := gomn.WrBin(db, dbPath); err != nil {
		eror(w, "failed to save db\n", err)
	};w.Write([]byte("saved db\n"))

	return nil
}

func dlBin(w http.ResponseWriter, typ string) {
	switch strings.ToLower(typ) {
	 case "bin", "b", "binary", "raw", "r", "gaas":
		file, err := os.Open(dbPath)
		if err != nil {
			eror(w, "opening db binary", err)
		}; defer file.Close()
		t := time.Now().Format("2006-01-02_15:04:05")
		sugFileName := fmt.Sprintf("binMap_db_%s.bgomn", t)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition",
			fmt.Sprintf("attachment; filename=\"%s\"", sugFileName))
		if _, err = io.Copy(w, file); err != nil {
			log.Error("err streaming binary to client")
		}
   case "key-value", "t", "text", "k-v", "kv", "key-val", "key value", "key_val", "pair", "p", "pairs":
		for key, val := range db {
			w.Write([]byte(fmt.Sprintf("%s = % x\n", key, val)))
		}
	 case "g", "gomn", "std", "standard", "s":
		for key, val := range db {
			w.Write([]byte(fmt.Sprintf("[\"%s\"] := \"% x\"\n", key, val)))
		}	
	default:
		log.Warnf("attempt to download db as unsupported type:  %s", typ) 
	}
	return
}
