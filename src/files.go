package main

/*
 * This is so I don't have to package extra files
 *  and for things that would be better to have in
 *   memory for speed
 */

import (
	"strconv" //used for default db 
	"github.com/Supraboy981322/gomn" //used for type
)

func defConfig() []byte {
	return []byte(`//configuration for binMap
["port"] := 4780
["log level"] := "debug"
["db path"] := "data/db.gaab"
["compress db"] := false`)
} 

func defDB() gomn.Map {
	return gomn.Map{
		"version": []byte("who knows"),
		"port": []byte(strconv.Itoa(port)),
		"foo.c": []byte(`#include <stdio.h>

int main(void) {
  printf("foo\n");
  return 0;
}`),
	}
}
