package limits

import (
	"io/ioutil"
	"log"
	"testing"

	"k8s.io/apimachinery/pkg/api/resource"
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

func TestBetweenCPU(t *testing.T) {
	lowqty, err := resource.ParseQuantity("1G")
	if err != nil {
		log.Fatal(err)
	}
	var highqty resource.Quantity
	highqty, err = resource.ParseQuantity("2G")
	if err != nil {
		log.Fatal(err)
	}
	var test resource.Quantity
	test, err = resource.ParseQuantity("1.5G")
	if err != nil {
		log.Fatal(err)
	}
	ls := LimitSentry{
		CPUMax: highqty,
		CPUMin: lowqty,
	}
	if !ls.BetweenCPU(test) {
		t.Error("expected qty for CPU to be in between")
	}

}
