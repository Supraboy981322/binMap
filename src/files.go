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

//valid: true or false (no quotes)
["log requests"] := true 

//set the path to your db
//  can be absolute or relative
["db path"] := "data/db.gaab"

//enable compression
//  true: use compression
//  false: no compression
["compress db"] := false

//set to 0 for never
["clear db every n seconds"] := 0 

//set to -1 for never
["clear db if size is n MB"] := -1

//when clearing the db,
//  should it be set to default?
["clear db to default"] := false

//disable to save RAM
["use in-memory db"] := true //disable to save RAM

//disable if you want lose 
//  db when server is stopped
["use disk db"] := true`)
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
