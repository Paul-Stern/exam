package main

import (
	"log"
	"os"
)

func clearTmp() (err error) {
	tmpfiles, err := os.ReadDir("tmp")
	if err != nil {
		return
	}
	for _, file := range tmpfiles {
		err = os.Remove(file.Name())
		if err != nil {
			return
		}
	}
	log.Println("clearTmp: tmp dir successfully cleared")
	return nil
}
