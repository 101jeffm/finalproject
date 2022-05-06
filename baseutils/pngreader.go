package baseutils

import (
	"bufio"
	"bytes"
	"log"
	"os"
)

//Processes image from file
func PNGPreProcess(d *os.File) (*bytes.Reader, error) {

	//Structure
	png, err := d.Stat()
	if err != nil {
		log.Fatal(err)
	}

	size := png.Size()
	buff := make([]byte, size)

	//Reader
	buffRead := bufio.NewReader(d)
	_, err = buffRead.Read(buff)
	if err != nil {
		log.Fatal(err)
	}
	buffReader := bytes.NewReader(buff)

	return buffReader, err
}
