package main

import (
  "bytes"
  // "encoding/binary"
  "io/ioutil"
  // "os"
  "fmt"
  "log"
  "bufio"
  // "io"
  // "encoding/binary"
)

var (
  ID3 = []byte{0x49,0x44,0x33}
  // TIT2 = []byte{54,49,54,32}
)

func main() {

  // var file string
  // fmt.Scanf("%s", &file)
  // fmt.Println(file)

  byt, err := ioutil.ReadFile("./cochise.mp3")
  if err != nil {
    log.Fatal(err)
  }

  b := bytes.NewReader(byt)

  br := bufio.NewReader(b)

  p := make([]byte, 3)
  a, _ := br.Read(p)
  fmt.Printf("% x  %i\n", p, a )

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


  // var dataLen int

  // for {

    // da, err := br.ReadBytes(0x54)
    // byts, err := br.Peek(3)
    // if err != nil {
    //   fmt.Println(err)
    // }
    //
    // if ok := bytes.Equal(byts, []byte{0x49,84,50}); ok {
    //   dataLen = len(da)
    //   break
    // }
  // }

//   fmt.Println(dataLen)
//   sizeOffset := dataLen+3
//   size := byt[sizeOffset:sizeOffset+4]
//   fmt.Printf("% x\n",size)
//
//   var size2 int32
//   buf := bytes.NewReader(size)
//   err2 := binary.Read(buf, binary.BigEndian, &size2)
//   if err2 != nil {
//     panic("asdas")
//   }
//
//   flags := byt[sizeOffset+4:sizeOffset+6]
//   fmt.Printf("% x\n", flags)
//
//   dataStart := sizeOffset+6
//   dataEnd := int32(dataStart)+size2
//   fmt.Printf("-%s-\n", byt[dataStart:dataEnd])
//
//   // var dataLen int
//   // for dataLen == 0 {
//   //   _, err := br.ReadBytes(0x54)
//   //   if err != nil {
//   //     fmt.Println(err)
//   //   }
//   //
//   //   byts, err := br.Peek(3)
//   //   if err != nil {
//   //     fmt.Println(err)
//   //   }
//   //
//   //   fmt.Println("%", byts)
//     
//     // break
//     // if ok := bytes.Equal(byts, []byte{49,54,32}); ok {
//     //   fmt.Println("Here")
//     // } else {
//     //   continue
//     // }
//     // break
//
//
//     // c, _ := br.ReadByte()
//     // if c != marker {
//     //   continue
//     // }
//     // dataLen = len(da)
//
}
