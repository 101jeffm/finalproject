package baseutils

import (
	"bytes"
	"finalproject/basemodels"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func WriteData(read *bytes.Reader, md *basemodels.CmdLineArgs, b []byte) {
	//Gets offset value
	fmt.Println(md.Offset)
	offset, err := strconv.ParseInt(md.Offset, 10, 64)
	fmt.Println(offset)
	if err != nil {
		log.Fatal(err)
	}

	//Create write file
	write, err := os.OpenFile(md.Output, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}

	//Go to beginning of png file
	read.Seek(0, 0)
	var buff = make([]byte, offset)
	//Read from old image and write to new image
	read.Read(buff)
	write.Write(buff)
	write.Write(b)
	if md.Decode {
		read.Seek(int64(len(b)), 1)
	}
	_, err = io.Copy(write, read)
	if err == nil {
		fmt.Printf("Success: %s created\n", md.Output)
	}
}
