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

func saveCert(data []byte) (name string, err error) {
	_, err = os.ReadDir("tmp")
	if os.IsNotExist(err) {
		os.Mkdir("tmp", 0755)
	}
	f, err := os.CreateTemp("./tmp", "cert-*.pdf")
	if err != nil {
		return
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		return
	}
	log.Printf("saveCert: successfully saved cert: %s\n", f.Name())
	return f.Name(), err
}
