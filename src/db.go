package main

import (
	"github.com/Supraboy981322/gomn"
)

func mapDB() error {
	var err error
	if db, err = gomn.ReadBin(dbPath); err != nil {
		return err
	}

	return nil
}

func updateDB(key string, value []byte) error {
	db[key] = value
	if err := gomn.WrBin(db, dbPath); err != nil {
		return err
	}
	return nil
}
