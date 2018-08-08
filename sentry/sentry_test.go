package sentry

import (
	"io/ioutil"
	"log"
	"testing"
)

var (
	testpod []byte
)

func init() {
	var err error
	testpod, err = ioutil.ReadFile("test.json")
	if err != nil {
		log.Fatal(err)
	}
}

func TestAdmit(t *testing.T) {

}
