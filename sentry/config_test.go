package sentry

import (
	"log"
	"testing"

	"github.com/davecgh/go-spew/spew"

	"k8s.io/apimachinery/pkg/api/resource"
)

func TestLoad(t *testing.T) {
	c := NewConfig()
	err := c.Load()
	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(c)
}

func TestBetween(t *testing.T) {
	m := MinMax{
		Min: "2G",
		Max: "5G",
	}
	err := m.SetResources()
	if err != nil {
		t.Fatal(err)
	}
	test, err := resource.ParseQuantity("1G")
	if err != nil {
		t.Fatal(err)
	}
	if m.Between(test) {
		log.Fatal("Expecting below value to fail")
	}
	test, err = resource.ParseQuantity("2G")
	if err != nil {
		t.Fatal(err)
	}
	if !m.Between(test) {
		log.Fatal("Expecting equal to Min value to pass")
	}
	test, err = resource.ParseQuantity("3G")
	if err != nil {
		t.Fatal(err)
	}
	if !m.Between(test) {
		log.Fatal("Expecting middle value to pass")
	}
	test, err = resource.ParseQuantity("5G")
	if err != nil {
		t.Fatal(err)
	}
	if !m.Between(test) {
		log.Fatal("Expecting equal to Max value to pass")
	}
	test, err = resource.ParseQuantity("6G")
	if err != nil {
		t.Fatal(err)
	}
	if m.Between(test) {
		log.Fatal("Expecting high value to fail")
	}
}
