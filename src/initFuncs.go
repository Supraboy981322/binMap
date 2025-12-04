package main

import (
	"os"
	"fmt"
	"errors"
	"strings"
	"path/filepath"
	"github.com/charmbracelet/log"
	"github.com/Supraboy981322/gomn"
)

func configure() error {
	var ok bool   //avoids golang quirk
	var err error //  (not sure if bug)

	//read the config
	if config, err = gomn.ParseFile(configPath); err != nil {
		return fmt.Errorf("parsing config:  %v", err)
	} else { log.Debug("parsed config") }

	//set the port
	if port, ok = config["port"].(int); !ok {
		return errors.New("assert port")
	} else { log.Debug("port set") }

	//get the log level 
	if logLvl, ok = config["log level"].(string); ok {
		log.Debug("read level")
		//set the log level
		switch strings.ToLower(logLvl) {
 		 case "debug": log.SetLevel(log.DebugLevel)
		 case "info":  log.SetLevel(log.InfoLevel)
		 case "warn":	 log.SetLevel(log.WarnLevel)
		 case "error": log.SetLevel(log.ErrorLevel)
		 case "fatal": log.SetLevel(log.FatalLevel)
		 default: log.SetLevel(log.DebugLevel) //default to debug
			//asume invalid, but don't exit
			log.Warn("invalid log level; defaulting to debug")
		};log.Infof("log level set to %s", log.GetLevel())
	} else { return errors.New("assert log level") }

	if _, ok := config["log requests"].(bool); !ok {
		return errors.New("asserting log request, "+
			"must be bool (true or false) with no quotes")
	}

	//errs silently, to (hopefully) prevent
	//  accidentally deleting db
	//    (which is the only extra permission granted)
	if foo, ok := config["admin ip"].([]any); ok {
		for _, adminRaw := range foo {
			if admin, ok := adminRaw.(string); ok {
				adminPermIP = append(adminPermIP, admin)
			}
		}
	}

	return nil
}

func initDB() error {
	var ok bool   //avoids golang quirk
	var err error //  (not sure if bug)

	//get the path of db from config
	if dbPath, ok = config["db path"].(string); ok {
		//ensure path exists 
		if err = os.MkdirAll(filepath.Dir(dbPath), 0777); err != nil {
			log.Errorf("MkdirAll:  %v", err)
			return fmt.Errorf("creating db path:  %v", err)
		} else { log.Debug("ensured db path exists") }

		//make sure there's no problems with db path
		if _, err = os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
			log.Warn("there appears to be no db file; creating one")

			//if no db, mk default db 
			m := defDB()
			if err = gomn.WrBin(m, dbPath); err != nil {
				log.Errorf("WrBin:  %v", err)
				return fmt.Errorf("writing default db to binary:  %v", err)
			} else { log.Debug("created default db") }
		} else { log.Debug("found db") }
		
		//read the db binary from disk
		if db, err = gomn.ReadBin(dbPath); err != nil {
			log.Errorf("ReadBin:  %v", err)
			return fmt.Errorf("reading db binary from disk:  %v", err)
		} else { log.Debug("read database") }
	} else { //return db path err
		log.Errorf("Stat(%s)", dbPath)
		return errors.New("assert db path from config")
	}

	return nil
}
