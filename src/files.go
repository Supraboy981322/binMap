package main

import "github.com/Supraboy981322/gomn"

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
		"foo.c": []byte(`#include <stdio.h>

int main(void) {
  printf("foo\n");
  return 0;
}`),
	}
}
