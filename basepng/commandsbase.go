package basepng

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"finalproject/baseutils"
	"fmt"
	"hash/crc32"
	"log"
	"os"
	"strconv"
	"strings"
)

//Header, Chunk, and Info structs
type PngHeader struct {
	Header uint64
}

type PngChunk struct {
	Size uint32
	Type uint32
	Data []byte
	CRC  uint32
}

type PngInfo struct {
	Chunk  PngChunk
	Offset int64
}

//Validate that header is png header
func (pi *PngInfo) ValidatePng(buff *bytes.Reader) {
	var head PngHeader

	if err := binary.Read(buff, binary.BigEndian, &head.Header); err != nil {
		log.Fatal(err)
	}

	//Makes array for header bytes
	buffArr := make([]byte, 8)
	binary.BigEndian.PutUint64(buffArr, head.Header)

	//Validates
	if string(buffArr[1:4]) != "PNG" {
		log.Fatal("Not PNG File")
	} else {
		fmt.Println("Valid PNG")
	}
}

//Returns offset without moving position in png
func (pi *PngInfo) GetOffset(b *bytes.Reader) {
	offset, _ := b.Seek(0, 1)
	pi.Offset = offset
}

func (pi *PngInfo) chunkTypeToString() string {
	s := fmt.Sprintf("%x", pi.Chunk.Type)
	decode, _ := hex.DecodeString(s)
	r := fmt.Sprintf("%s", decode)
	return r
}

//4 functions that read chunk info from the chunk
func (pi *PngInfo) ChunkSizeR(buff *bytes.Reader) {
	if err := binary.Read(buff, binary.BigEndian, &pi.Chunk.Size); err != nil {
		log.Fatal(err)
	}
}

func (pi *PngInfo) ChunkTypeR(buff *bytes.Reader) {
	if err := binary.Read(buff, binary.BigEndian, &pi.Chunk.Type); err != nil {
		log.Fatal(err)
	}
}

func (pi *PngInfo) ChunkDataR(buff *bytes.Reader, chunkLen uint32) {
	pi.Chunk.Data = make([]byte, chunkLen)
	if err := binary.Read(buff, binary.BigEndian, &pi.Chunk.Data); err != nil {
		log.Fatal(err)
	}
}

func (pi *PngInfo) ChunkCRCR(buff *bytes.Reader) {
	if err := binary.Read(buff, binary.BigEndian, &pi.Chunk.CRC); err != nil {
		log.Fatal(err)
	}
}

func (pi *PngInfo) ChunkRead(buff *bytes.Reader) {
	pi.ChunkSizeR(buff)
	pi.ChunkTypeR(buff)
	pi.ChunkDataR(buff, pi.Chunk.Size)
	pi.ChunkCRCR(buff)
}

//Turn Type from string to int
func (pi *PngInfo) StrToInt(s string) uint32 {
	si := []byte(s)
	return binary.BigEndian.Uint32(si)
}

//Assign chunk size
func (pi *PngInfo) CreateSize() uint32 {
	return uint32(len(pi.Chunk.Data))
}

//Calculates checksum on bytes
func (pi *PngInfo) CreateCRC() uint32 {
	by := new(bytes.Buffer)
	if err := binary.Write(by, binary.BigEndian, pi.Chunk.Type); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(by, binary.BigEndian, pi.Chunk.Data); err != nil {
		log.Fatal(err)
	}
	return crc32.ChecksumIEEE(by.Bytes())
}

//Get all chunk segment data into single buffer
func (pi *PngInfo) Marshal() *bytes.Buffer {
	by := new(bytes.Buffer)
	if err := binary.Write(by, binary.BigEndian, pi.Chunk.Size); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(by, binary.BigEndian, pi.Chunk.Type); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(by, binary.BigEndian, pi.Chunk.Data); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(by, binary.BigEndian, pi.Chunk.CRC); err != nil {
		log.Fatal(err)
	}
	return by
}

//Find Ancillary Chunks
func (pi *PngInfo) checkCritType() string {
	fChar := string([]rune(pi.chunkTypeToString())[0])
	if fChar == strings.ToUpper(fChar) {
		return "Critical"
	}
	return "Ancillary"
}

//See if Offset if valid
func isValidOff(s string, ancil []string) bool {
	for _, ss := range ancil {
		if ss == s {
			return true
		}
	}
	return false
}

//Process bytes after png header
func (pi *PngInfo) PngProcess(buff *bytes.Reader) {
	//Start with Validation
	pi.ValidatePng(buff)

	var code string
	var key string
	var message string
	var offset string
	var ty string
	ancil := make([]string, 0)

	scanner := bufio.NewReader(os.Stdin)

	count := 1
	chunkType := ""
	endChunkType := "IEND" //last type represented before eof for png

	fmt.Println("--------------------------------------------------")
	fmt.Print("Would you like to encode or decode this image(E/D)?: ")
	fmt.Scanln(&code)
	if code != "E" && code != "D" {
		fmt.Println("Invalid entry, default to Encode")
		code = "E"
	}

	fmt.Println("--------------------------------------------------")
	fmt.Print("Xor or Aes Method(X/A)?: ")
	fmt.Scanln(&ty)
	if ty != "X" && ty != "A" {
		fmt.Println("Invalid entry, default to Xor")
		ty = "X"
	}

	fmt.Println("--------------------------------------------------")
	fmt.Print("Please enter your message key: ")
	fmt.Scanln(&key)
	if key == "" {
		fmt.Println("Invalid entry, default to key")
		key = "key"
	}

	if code == "E" { //If encoding
		var p PngInfo
		for chunkType != endChunkType {
			offset = "0x"
			pi.GetOffset(buff)
			pi.ChunkRead(buff)

			//Not picking correct chunks...
			//For splitting into several bytes -- want to take first 2 Ancillary chunks to store
			//if pi.checkCritType() == "Ancillary" {
			//	offset = offset + strconv.FormatInt(pi.Offset, 16)
			//	ancil = append(ancil, offset)
			//}
			//Defaulting to chunk IEOF since that works...
			if pi.chunkTypeToString() == endChunkType {
				offset = offset + strconv.FormatInt(pi.Offset, 16)
				ancil = append(ancil, offset)
			}
			chunkType = pi.chunkTypeToString()
			count++
		}

		fmt.Println("--------------------------------------------------")
		if len(ancil) == 1 {
			fmt.Println("Offset defaulted to: " + ancil[0])
			offset = ancil[0]
		} else {
			fmt.Println(strings.Join(ancil, ", "))
			fmt.Print("Pick an offset from the list above: ")
			fmt.Scanln(&offset)
			//for !isValidOff(offset, ancil) {
			//	fmt.Print("Invalid choice, pick again: ")
			//	fmt.Scanln(&offset)
			//}

		}

		byteOffset, _ := strconv.ParseInt(offset, 0, 64)
		offset = strconv.FormatInt(byteOffset, 10)
		fmt.Println("--------------------------------------------------")
		fmt.Print("Please enter your message to encode: ")
		//can now read entire line for message
		message, _ = scanner.ReadString('\n')
		if message == "" {
			log.Fatal("Cannot encode.")
		} else if ty == "A" && len(message) < 16 {
			log.Fatal("Cannot encode.")
		}

		p.Chunk.Data = baseutils.EncoderDecoder([]byte(message), key, ty, code)
		p.Chunk.Type = p.StrToInt("rNDm")
		p.Chunk.Size = p.CreateSize()
		p.Chunk.CRC = p.CreateCRC()
		buffp := p.Marshal()
		buffpp := buffp.Bytes()
		fmt.Println("--------------------------------------------------")
		fmt.Printf("Payload Original: % X\n", []byte(message))
		fmt.Printf("Payload Encode: % X\n", p.Chunk.Data)
		fmt.Println("--------------------------------------------------")
		baseutils.WriteData(buff, offset, false, buffpp)
	} else { //If Decoding
		var p PngInfo
		fmt.Println("--------------------------------------------------")
		fmt.Print("Please enter your offset: ")
		fmt.Scanln(&offset)

		if offset == "" {
			fmt.Println("Invalid offset, default to 0x85258")
			offset = "0x85258"
		}

		byteOffset, _ := strconv.ParseInt(offset, 0, 64)
		offset = strconv.FormatInt(byteOffset, 10)

		off, _ := strconv.ParseInt(offset, 10, 64)
		buff.Seek(off, 0)
		p.ChunkRead(buff)
		ogData := p.Chunk.Data
		p.Chunk.Data = baseutils.EncoderDecoder(p.Chunk.Data, key, ty, code)
		p.Chunk.CRC = p.CreateCRC()
		buffp := p.Marshal()
		buffpp := buffp.Bytes()
		fmt.Println("--------------------------------------------------")
		fmt.Printf("Payload Original: % X\n", ogData)
		fmt.Printf("Payload Decode: % X\n", p.Chunk.Data)
		fmt.Printf("Original Message: %s\n", string(p.Chunk.Data))
		fmt.Println("--------------------------------------------------")

		baseutils.WriteData(buff, offset, true, buffpp)
	}
}
