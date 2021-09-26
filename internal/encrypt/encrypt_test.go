package encrypt

import "testing"

func TestEncoding(t *testing.T) {
	//paranoia
	scaryKey = ""
	err := SetKey("hello")
	if err != nil {
		t.Fatal("set key failed", err)
	}
	if scaryKey != "hello" {
		t.Fatal("bad key setting", scaryKey)
	}

	err = SetKey("bye")
	if err == nil || scaryKey != "bye" {
		t.Fatal("Resetting Key didnt work as expected")
	}

	source := "encoding is fun for :: people"
	out := Encode(source)
	if out == source {
		t.Fatal("Um encoding should do something...")
	}

	out2, err := Decode(out)
	if err != nil {
		t.Fatal("Decoding failed", err)
	}
	if out2 != source {
		t.Fatal("encode decode failed")
	}
}
