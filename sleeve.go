package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jamesduncombe/colour"
	"github.com/jamesduncombe/sleeve/cover"
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
	showDump, verbose     bool
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

func parseFlags() {
	flag.Parse()

	if inputFile == "" || outputFile == "" {
		colour.Red(header)
		flag.PrintDefaults()
		os.Exit(0)
	}
}

func main() {

	parseFlags()

	f, _ := openFile(inputFile)
	defer f.Close()

	// we want to exit very early if this isn't actually an ID3 holding MP3

	p := make([]byte, 3)
	f.Read(p)

	if !cover.HasId3(p) {
		msg("Can't find ID3")
		return
	}

	// ok, good to go, read it all and setup buffers and readers
	byt, _ := ioutil.ReadAll(f)
	b := bytes.NewReader(byt)
	br := bufio.NewReader(b)

	// get version of ID3
	if showDump || verbose {
		cover.Id3Ver(br)
	}

	if !cover.HasPicture(br) {
		msg("Can't find picture")
		return
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

func openFile(filePath string) (*os.File, error) {
	msg("Input file: " + filePath)

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// make sure this isn't a directory
	d, _ := os.Stat(filePath)
	if d.IsDir() {
		log.Fatal("This is a directory, please pass a file instead")
	}

	return f, err
}

func msg(message string) {
	if verbose {
		fmt.Println(message)
	}
}
