package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

const (
	payloadWithArray  = `[{"stuff":"something"},{"stuff":"something else"}]`
	payloadWithObject = `{"stuff":"one thing"}`
)

type stuff struct {
	Stuff string
}

func main() {
	bodies := []io.Reader{
		strings.NewReader(payloadWithArray),
		strings.NewReader(payloadWithObject),
	}

	stuffsToReturn := []stuff{}
	// level := 0
	for _, body := range bodies {
		// reads everything to memory..
		// if payload is too big you can use teereader etc.
		// to do it as a stream
		payload, err := ioutil.ReadAll(body)
		if err != nil {
			log.Fatal(err)
		}

		buf := bytes.NewBuffer(payload)

		// peek into the payload through a decoder
		dec := json.NewDecoder(buf)
		t, err := dec.Token()
		if err != nil {
			log.Fatal(err)
		}

		if d, isDelim := t.(json.Delim); isDelim {
			switch d {
			case '[':
				newStuffs := []stuff{}
				err := json.Unmarshal(payload, &newStuffs)
				if err != nil {
					log.Fatal(err)
				}
				stuffsToReturn = append(stuffsToReturn, newStuffs...)
			case '{':
				newStuff := stuff{}
				err := json.Unmarshal(payload, &newStuff)
				if err != nil {
					log.Fatal(err)
				}
				stuffsToReturn = append(stuffsToReturn, newStuff)
			}
		}
	}
	fmt.Println(stuffsToReturn)
}
