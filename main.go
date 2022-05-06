package main

import (
	//"finalproject/basemodels"
	"finalproject/basepng"
	"finalproject/baseutils"
	"fmt"
	"log"
	"os"
	//"strconv"
)

var (
	//opts  basemodels.CmdLineArgs
	png basepng.PngInfo
)

func main() {
	var input string
	fmt.Println("--------------------------------------------------")
	fmt.Println("Welcome to the Encoder/Decoder")
	fmt.Print("Enter image location: ")
	fmt.Scanln(&input)
	dat, err := os.Open(input)
	defer dat.Close()
	bReader, err := baseutils.PNGPreProcess(dat)
	if err != nil {
		log.Fatal(err)
	}
	png.PngProcess(bReader)
}
