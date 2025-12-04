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
	return []byte(`//enter your port number here
["port"] := 4780

//enter your log level
//  valid options:
//    - "debug"
//    - "info" //recommended
//    - "warn"
//    - "error"
//    - "fatal"
["log level"] := "debug"
["log requests"] := true //valid: true or false (no quotes)

//set the path to your db
//  can be absolute or relative
["db path"] := "data/db.gaab"

//enable compression
//  true: use compression
//  false: no compression
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
