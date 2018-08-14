package limits

import (
	"io/ioutil"
	"log"
)

var (
	testpod []byte
)

func init() {
	var err error
	testpod, err = ioutil.ReadFile("podtest.json")
	if err != nil {
		log.Fatal(err)
	}
}
