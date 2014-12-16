package main

import (
  "bufio"
  "bytes"
  "flag"
  "fmt"
  "io/ioutil"
  "log"
  "os"

  "github.com/jamesduncombe/sleeve/cover"
  "github.com/jamesduncombe/colour"
)

const header = `
         ______                      
  __________  /____________   ______ 
  __  ___/_  /_  _ \  _ \_ | / /  _ \
  _(__  )_  / /  __/  __/_ |/ //  __/
  /____/ /_/  \___/\___/_____/ \___/ 
`

var (
  inputFile, outputFile string
  showDump, verbose bool
)

var (
  JPEG_END = []byte{0xff, 0xd9}
)

func init() {
  flag.StringVar(&inputFile, "i", "", "Input filename of the MP3 to read (or Stdin)")
  flag.StringVar(&outputFile, "o", "", "Output filename of the JPEG")
  flag.BoolVar(&verbose, "v", false, "Verbose output")
  flag.BoolVar(&showDump, "d", false, "Show dump of ID3 info")
}

func main() {

  flag.Parse()

  if inputFile == "" || outputFile =="" {
    colour.Red(header)
    flag.PrintDefaults()
    os.Exit(0)
  }

  // open file
  msg("Input file: " + inputFile)
  byt, _ := openFile(inputFile)

  // setup buffers and readers
  b := bytes.NewReader(byt)
  br := bufio.NewReader(b)

  // read initial 3 bytes
  p := make([]byte, 3)
  br.Read(p)

  // check for ID3
  if !cover.HasId3(p) {
    msg("Can't find ID3")
    return
  }

  // get version of ID3
  if showDump || verbose {
    cover.Id3Ver(br)
  }

  if ok := cover.HasPicture(br); !ok {
    msg("Can't find picture")
    os.Exit(1)
  }

  jpegData := bytes.NewBuffer([]byte{0xff}) // JPEG header SOI

  msg("Looking for JPEG")
  for {
    br.ReadBytes(0xff)
    by, _ := br.Peek(2)
    if ok := bytes.Equal(by, []byte{0xd8, 0xff}); ok {
      msg("Found JPEG data")
      for {
        peaky, _ := br.Peek(2)
        if ok := bytes.Equal(peaky, JPEG_END); ok {
          msg("Saving output file: " + outputFile)
          f, err := os.Create(outputFile)
          if err != nil {
            fmt.Println(err)
          }
          defer f.Close()
          jpegData.Write(JPEG_END)
          f.Write(jpegData.Bytes())

          break
        }
        c, err := br.ReadByte()
        if err != nil {
          panic("Can't progress")
        }
        jpegData.WriteByte(c)
      }
      break
    } else {
      break
    }
  }

}

func openFile(filePath string) ([]byte, error) {

  byt, err := ioutil.ReadFile(filePath)
  if err != nil {
    log.Fatal(err)
  }
  return byt, err
}

func msg(message string) {
  if verbose {
    fmt.Println(message)
  }
}
