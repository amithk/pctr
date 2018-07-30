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
