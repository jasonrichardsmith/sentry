package limits

import (
	"testing"

	"github.com/jasonrichardsmith/sentry/sentry"
)

func TestLoadSentry(t *testing.T) {
	c := Config{
		Memory: MinMax{
			Min: "1g",
			Max: "1g",
		},
		CPU: MinMax{
			Min: "1g",
			Max: "1g",
		},
	}
	_, err := c.LoadSentry()
	if err == nil {
		t.Fatal("Expecting resource parse error for Memory Max")
	}
	c.Memory.Max = "1G"
	_, err = c.LoadSentry()
	if err == nil {
		t.Fatal("Expecting resource parse error Memory Min")
	}
	c.Memory.Min = "1G"
	_, err = c.LoadSentry()
	if err == nil {
		t.Fatal("Expecting resource parse error CPU Min")
	}
	c.CPU.Min = "1G"
	_, err = c.LoadSentry()
	if err == nil {
		t.Fatal("Expecting resource parse error CPU Max")
	}
	c.CPU.Max = "1G"
	var s sentry.Sentry
	s, err = c.LoadSentry()
	ls := s.(LimitSentry)
	if err != nil {
		t.Fatal(err)
	}
	if ls.CPUMax.String() != c.CPU.Max {
		t.Fatal("CPU Max mismatch")
	}
	if ls.CPUMin.String() != c.CPU.Min {
		t.Fatal("CPU Min mismatch")
	}
	if ls.MemoryMax.String() != c.Memory.Max {
		t.Fatal("Memory Max mismatch")
	}
	if ls.MemoryMin.String() != c.Memory.Min {
		t.Fatal("Memory Min mismatch")
	}

}
