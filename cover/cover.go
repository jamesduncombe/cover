// Package cover implements utilities for dealing with ID3 tags in MP3s,
// it provides helper methods for this.
package cover

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

var (
	ID3  = []byte{0x49, 0x44, 0x33}
	APIC = []byte{0x41, 0x50, 0x49, 0x43}
)

type Cover struct {
}

// HasId3 returns whether the passed in bytestream contains what is an ID3 tag
func HasId3(p []byte) bool {
	return bytes.Equal(p, ID3)
}

func Id3Ver(br *bufio.Reader) {
	ver := make([]byte, 2)
	br.Read(ver)
	fmt.Printf("ID3v2.%d Rev: %d\n", uint8(ver[0]), uint8(ver[1]))
}

func HasPicture(br *bufio.Reader) bool {
	// find picture
	for {
		_, err := br.ReadBytes(APIC[0])
		// if we hit the end then it didn't have a APIC tag :(
		if err == io.EOF {
			return false
		}
		// otherwise, peek 3 ahead
		by, _ := br.Peek(3)
		if ok := bytes.Equal(by, APIC[1:4]); ok {
			return true
		}
	}
}
