package main

import (
  "bufio"
  "bytes"
  "fmt"
  "io"
  "io/ioutil"
  "log"
  "os"
)

var (
  ID3      = []byte{0x49, 0x44, 0x33}
  APIC     = []byte{0x41, 0x50, 0x49, 0x43}
  JPEG_END = []byte{0xff, 0xd9}
)

func openFile(filePath string) ([]byte, error) {

  byt, err := ioutil.ReadFile(filePath)
  if err != nil {
    log.Fatal(err)
  }
  return byt, err
}

func main() {

  // open file
  filePath := os.Args[1]
  byt, _ := openFile(filePath)

  // setup buffers and readers
  b := bytes.NewReader(byt)
  br := bufio.NewReader(b)

  // read initial 3 bytes
  p := make([]byte, 3)
  br.Read(p)

  // check for ID3
  if ok := bytes.Equal(p, ID3); ok {
    fmt.Println("Is ID3")
  } else {
    fmt.Println("Not ID3")
    return
  }

  // get version of ID3
  ver, _ := br.ReadByte()
  fmt.Printf("ID3v2.%d\n", uint8(ver))

  // get revision
  revision, _ := br.ReadByte()
  fmt.Printf("Revision: %d\n", uint8(revision))

  // get flags
  flags, _ := br.ReadByte()
  fmt.Printf("Flags: %x \n", flags)

  // find picture
  for {
    _, err := br.ReadBytes(APIC[0])
    if err == io.EOF {
      fmt.Println("EOF reached")
      return
    }
    by, _ := br.Peek(3)
    if ok := bytes.Equal(by, APIC[1:4]); ok {
      fmt.Println("Found APIC Tag")
      break
    }
  }

  jpegData := bytes.NewBuffer([]byte{0xff}) // JPEG header SOI

  // find mime type
  for {
    br.ReadBytes(0xff)
    by, _ := br.Peek(2)
    if ok := bytes.Equal(by, []byte{0xd8, 0xff}); ok {
      fmt.Println("Found JPEG")
      for {
        peaky, _ := br.Peek(2)
        if ok := bytes.Equal(peaky, JPEG_END); ok {
          fmt.Println("End of JPEG")

          f, err := os.Create("./" + filePath + ".jpg")
          if err != nil {
            fmt.Println(err)
          }
          defer f.Close()
          fmt.Println("Outputting JPEG")
          jpegData.Write(JPEG_END)
          f.Write(jpegData.Bytes())

          break
        }
        c, err := br.ReadByte()
        if err != nil {
          fmt.Println("Woops")
        }
        jpegData.WriteByte(c)
      }
      break
    } else {
      break
    }
  }

}
