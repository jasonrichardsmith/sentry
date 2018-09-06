package limits

import (
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/api/resource"
)

func TestName(t *testing.T) {
	c := Config{}
	if c.Name() != "limits" {
		t.Fatal("Failed name test")
	}
}

func TestQtyHook(t *testing.T) {
	var i int
	test := "1G"
	vtype := reflect.TypeOf(i)
	v, err := QtyHookFunc(vtype, vtype, test)
	if err != nil {
		t.Fatal(err)
	}
	if v.(string) != test {
		t.Fatal("Expecting skip on not being string")
	}
	stype := reflect.TypeOf(test)
	v, err = QtyHookFunc(stype, vtype, test)
	if err != nil {
		t.Fatal(err)
	}
	if v.(string) != test {
		t.Fatal("Expecting skip on not being qty")
	}
	vtype = reflect.TypeOf(resource.Quantity{})
	v, err = QtyHookFunc(stype, vtype, test)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.TypeOf(resource.Quantity{}) != reflect.TypeOf(v) {
		t.Fatal("Expecting return of qty")
	}
}
