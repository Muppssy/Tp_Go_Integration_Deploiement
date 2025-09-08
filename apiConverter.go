package main

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

func yml2json() {

	yfile, err := ioutil.ReadFile("openapi.yml")

	if err != nil {
		log.Fatal(err)
	}

	var data any

	err = yaml.Unmarshal(yfile, &data)

	if err != nil {
		log.Fatal(err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "\t")
	enc.Encode(data)
}
