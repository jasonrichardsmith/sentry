package limits

import (
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	podpass    []byte
	podhighcpu []byte
	podhighmem []byte
	podnocpu   []byte
	podnomem   []byte
	lowqty     resource.Quantity
	highqty    resource.Quantity
	test       resource.Quantity
)

func init() {
	var err error
	podpass, err = ioutil.ReadFile("podtest.json.pass")
	if err != nil {
		log.Fatal(err)
	}
	podhighcpu, err = ioutil.ReadFile("podtest.json.highcpu")
	if err != nil {
		log.Fatal(err)
	}
	podhighmem, err = ioutil.ReadFile("podtest.json.highmem")
	if err != nil {
		log.Fatal(err)
	}
	podnocpu, err = ioutil.ReadFile("podtest.json.nocpu")
	if err != nil {
		log.Fatal(err)
	}
	podnomem, err = ioutil.ReadFile("podtest.json.nomem")
	if err != nil {
		log.Fatal(err)
	}
	lowqty, err = resource.ParseQuantity("1G")
	if err != nil {
		log.Fatal(err)
	}
	highqty, err = resource.ParseQuantity("2G")
	if err != nil {
		log.Fatal(err)
	}
	test, err = resource.ParseQuantity("1.5G")
	if err != nil {
		log.Fatal(err)
	}
}

func TestType(t *testing.T) {
	is := LimitSentry{}
	if is.Type() != "Pod" {
		t.Fatal("Failed type test")
	}
}

func TestBetweenCPU(t *testing.T) {
	c := Config{
		CPU: MinMax{
			Max: highqty,
			Min: lowqty,
		},
	}
	ls := c.LoadSentry().(LimitSentry)
	if !ls.BetweenCPU(test) {
		t.Error("expected qty for CPU to be in between")
	}
	ls.CPU.Max = lowqty
	if ls.BetweenCPU(test) {
		t.Error("expected qty for CPU to not be in between")
	}

}

func TestBetweenMemory(t *testing.T) {
	c := Config{
		Memory: MinMax{
			Max: highqty,
			Min: lowqty,
		},
	}
	ls := c.LoadSentry().(LimitSentry)
	if !ls.BetweenMemory(test) {
		t.Error("expected qty for Memory to be in between")
	}
	ls.Memory.Max = lowqty
	if ls.BetweenMemory(test) {
		t.Error("expected qty for Memory to not be in between")
	}

}

func TestAdmit(t *testing.T) {
	c := Config{
		Memory: MinMax{
			Max: highqty,
			Min: lowqty,
		},
		CPU: MinMax{
			Max: highqty,
			Min: lowqty,
		},
	}
	ls := c.LoadSentry().(LimitSentry)
	ar := v1beta1.AdmissionReview{
		Request: &v1beta1.AdmissionRequest{
			Object: runtime.RawExtension{
				Raw: podpass,
			},
		},
	}
	resp := ls.Admit(ar)
	if !resp.Allowed {
		t.Fatal("expected passing review")
	}
	ar.Request.Object.Raw = podhighmem
	resp = ls.Admit(ar)
	if resp.Allowed {
		t.Fatal("Expected highmem to fail")
	}
	if resp.Result.Message != LimitsOutsideMemory {
		t.Fatal("Expected memory out of range error message")
	}
	ar.Request.Object.Raw = podhighcpu
	resp = ls.Admit(ar)
	if resp.Allowed {
		t.Fatal("Expected highpod to fail")
	}
	if resp.Result.Message != LimitsOutsideCPU {
		t.Fatal("Expected cpu out of range error message")
	}
	ar.Request.Object.Raw = podnocpu
	resp = ls.Admit(ar)
	if resp.Allowed {
		t.Fatal("Expected limits not set error")
	}
	if resp.Result.Message != LimitsNotPresent {
		t.Fatal("Expected limits not set error message")
	}
	ar.Request.Object.Raw = podnomem
	resp = ls.Admit(ar)
	if resp.Allowed {
		t.Fatal("Expected limits not set error")
	}
	if resp.Result.Message != LimitsNotPresent {
		t.Fatal("Expected limits not set error message")
	}
	ar.Request.Object.Raw = podpass[0:5]
	resp = ls.Admit(ar)
	if !strings.Contains(resp.Result.Message, "json parse error") {
		t.Fatal("Expecting json parse error")
	}
}
