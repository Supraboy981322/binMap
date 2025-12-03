package main

import (
	"os"
	"errors"
	"strings"
	"path/filepath"
	"github.com/charmbracelet/log"
	"github.com/Supraboy981322/gomn"
)

func configure() error {
	var ok bool
	var err error

	if config, err = gomn.ParseFile(configPath); err != nil {
		return err
	} else { log.Debug("parsed config") }

	if port, ok = config["port"].(int); !ok {
		return errors.New("assert port")
	} else { log.Debug("port set") }

	if logLvl, ok = config["log level"].(string); ok {
		log.Debug("read level")
		switch strings.ToLower(logLvl) {
 		 case "debug": log.SetLevel(log.DebugLevel)
		 case "info":  log.SetLevel(log.InfoLevel)
		 case "warn":	 log.SetLevel(log.WarnLevel)
		 case "error": log.SetLevel(log.ErrorLevel)
		 case "fatal": log.SetLevel(log.FatalLevel)
		 default: log.SetLevel(log.DebugLevel)
			log.Warn("invalid log level; defaulting to debug")
		};log.Info("log level set")
	} else { return errors.New("assert log level") }

	if dbPath, ok = config["db path"].(string); ok {
		if err = os.MkdirAll(filepath.Dir(dbPath), 0777); err != nil {
			return err
		} else { log.Debug("ensured db path exists") }

		if _, err = os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
			log.Warn("there appears to be no db file; creating one")

			m := gomn.Map{"version": []byte("who knows")}
			if err = gomn.WrBin(m, dbPath); err != nil {
				return err
			} else { log.Debug("created default db") }
		} else { log.Debug("found db") }
		
		if db, err = gomn.ReadBin(dbPath); err != nil {
			return err
		} else { log.Debug("read database") }
	} else { return errors.New("assert db path") }

	return nil
}
