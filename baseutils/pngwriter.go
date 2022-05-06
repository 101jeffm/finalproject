package baseutils

import (
	"bytes"
	//"finalproject/basemodels"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func WriteData(read *bytes.Reader, off string, decode bool, b []byte) {
	var save string

	//Gets offset value
	offset, err := strconv.ParseInt(off, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Enter filename to save to: ")
	fmt.Scanln(&save)
	if save == "" && decode == true {
		fmt.Println("Invalid entry, default to decode.png")
		save = "decode.png"
	} else if save == "" && decode != true {
		fmt.Println("Invalid entry, default to encode.png")
		save = "encode.png"
	}
	//Create write file
	write, err := os.OpenFile(save, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}

	//Go to beginning of png file
	read.Seek(0, 0)
	var buff = make([]byte, offset)

	//Read from old image and write to new image
	read.Read(buff)
	write.Write(buff)

	//left shift so we overwrite instead of insert data
	if !decode {
		read.Seek(0-int64(len(b)), 1)
	}
	write.Write(b)
	if decode {
		read.Seek(int64(len(b)), 1)
	}
	_, err = io.Copy(write, read)
	if err == nil {
		fmt.Printf("Success: %s created\n", save)
	}
}
