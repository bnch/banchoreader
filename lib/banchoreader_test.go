package banchoreader

import (
	"os"
	"testing"
)

const testPackets = "\x03\x00\x00\x04\x00\x00\x00\x05\x0f\x00\x00" +
	"\x04\xa0\x00\x05\x00\x00\x00\x05\x0f\x00\x00\x06" +
	"\x0a\x00\x00\x01\x00\x00\x00\xad" +
	"\x0b\x00\x00\x04\x00\x00\x00\xde\xad\xbe\xef" +
	"\x0f\x00\x00\x04\x00\x00\x00\x32\x30\x3f\x4a"
const mustFail = "This shall not succeed"

func TestRead(t *testing.T) {
	b := []byte(testPackets)

	d := New()
	err := d.Dump(os.Stdout, b)
	if err != nil {
		t.Fatal(err)
	}

	d.Colored = true
	err = d.Dump(os.Stdout, b)
	if err != nil {
		t.Fatal(err)
	}

	d.Colored = false
	d.Ignored = []int{3, 10}
	err = d.Dump(os.Stdout, b)
	if err != nil {
		t.Fatal(err)
	}

	d.Ignored = []int{}
	d.IndentationSize = 4
	err = d.Dump(os.Stdout, b)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMust(t *testing.T) {
	defer func() {
		c := recover()
		if c == nil {
			t.Fatal("The dump succeded :(")
		}
	}()
	d := New()
	d.Colored = true
	d.MustDump(os.Stdout, []byte(mustFail))
}
