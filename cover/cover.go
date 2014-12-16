package cover

import (
  "bytes"
  "bufio"
  "fmt"
  "io"
)

var (
  ID3  = []byte{ 0x49, 0x44, 0x33 }
  APIC = []byte{ 0x41, 0x50, 0x49, 0x43 }
)

func HasId3(p []byte) bool {
  if ok := bytes.Equal(p, ID3); ok {
    return true
  } else {
    return false
  }
}

func Id3Ver(br *bufio.Reader) {
  ver := make([]byte, 2)
  br.Read(ver)
  fmt.Printf("ID3v2.%d Rev: %d\n", uint8(ver[0]), uint8(ver[1]))
}

func Flags() {
  // flags, _ := br.ReadByte()
  // fmt.Printf("Flags: %x \n", flags)
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
