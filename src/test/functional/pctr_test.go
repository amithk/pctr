package pctrtest

import "fmt"
import "os"
import "pctr"
import "testing"

func TestBasicIncrFunctionality(t *testing.T) {
	os.RemoveAll("/tmp/1234")
	pc, err := pctr.NewPersistentCounter("/tmp", "1234")
	if err != nil {
		t.Fatalf("error %v during NewPersistentCounter", err)
	}

	val, err1 := pc.IncrementValue(10)
	if err1 != nil {
		t.Fatalf("error %v during IncrementValue", err1)
	}

	if val != 10 {
		t.Fatalf("Unexpected counter value")
	}

	rval, err2 := pc.GetValue()
	if err2 != nil {
		t.Fatalf("error %v during GetValue", err2)
	}

	if rval != val {
		t.Fatalf("Unexpected counter value in file")
	}

	fmt.Println(rval)
}

func TestBasicGetNextFunctionality(t *testing.T) {
	os.RemoveAll("/tmp/12345")
	pc, err := pctr.NewPersistentCounter("/tmp", "12345")
	if err != nil {
		t.Fatalf("error %v during NewPersistentCounter", err)
	}

	for i := 0; i < 70; i++ {
		nv, err1 := pc.GetNext()
		if err1 != nil {
			t.Fatalf("Unexpected error %v during GetNext", err1)
		}
		if uint64(i+1) != nv {
			t.Fatalf("Unexpected value %v, %v in GetNext", i, nv)
		}
	}
}
