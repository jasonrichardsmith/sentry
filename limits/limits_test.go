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

func TestBetweenCPU(t *testing.T) {
	ls := LimitSentry{
		CPUMax: highqty,
		CPUMin: lowqty,
	}
	if !ls.BetweenCPU(test) {
		t.Error("expected qty for CPU to be in between")
	}
	ls.CPUMax = lowqty
	if ls.BetweenCPU(test) {
		t.Error("expected qty for CPU to not be in between")
	}

}

func TestBetweenMemory(t *testing.T) {
	ls := LimitSentry{
		MemoryMax: highqty,
		MemoryMin: lowqty,
	}
	if !ls.BetweenMemory(test) {
		t.Error("expected qty for Memory to be in between")
	}
	ls.MemoryMax = lowqty
	if ls.BetweenMemory(test) {
		t.Error("expected qty for Memory to not be in between")
	}

}

func TestAdmit(t *testing.T) {
	ls := LimitSentry{
		MemoryMax: highqty,
		MemoryMin: lowqty,
		CPUMax:    highqty,
		CPUMin:    lowqty,
	}
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
