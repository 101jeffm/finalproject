package basepng

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"finalproject/basemodels"
	"finalproject/baseutils"
	"fmt"
	"hash/crc32"
	"log"
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

//Not Needed, Remove later
func (pi *PngInfo) checkCritType() string {
	fChar := string([]rune(pi.chunkTypeToString())[0])
	if fChar == strings.ToUpper(fChar) {
		return "Critical"
	}
	return "Ancillary"
}

//Process bytes after png header
func (pi *PngInfo) PngProcess(buff *bytes.Reader, md *basemodels.CmdLineArgs) {
	//Start with Validation
	pi.ValidatePng(buff)

	if (md.Offset != "") && (md.Encode == false && md.Decode == false) {
		//not encoding or decoding
		var p PngInfo
		p.Chunk.Data = []byte(md.Payload)
		p.Chunk.Type = p.StrToInt(md.Type)
		p.Chunk.Size = p.CreateSize()
		p.Chunk.CRC = p.CreateCRC()
		buffp := p.Marshal()
		buffpp := buffp.Bytes()
		fmt.Printf("Payload Original: % X\n", []byte(md.Payload))
		fmt.Printf("Payload : % X\n", p.Chunk.Data)
		baseutils.WriteData(buff, md, buffpp)
	}
	if (md.Offset != "") && md.Encode { //If encoding
		var p PngInfo
		p.Chunk.Data = baseutils.EncoderDecoder([]byte(md.Payload), md.Key)
		p.Chunk.Type = p.StrToInt(md.Type)
		p.Chunk.Size = p.CreateSize()
		p.Chunk.CRC = p.CreateCRC()
		buffp := p.Marshal()
		buffpp := buffp.Bytes()
		fmt.Printf("Payload Original: % X\n", []byte(md.Payload))
		fmt.Printf("Payload Encode: % X\n", p.Chunk.Data)
		baseutils.WriteData(buff, md, buffpp)
	}
	if (md.Offset != "") && md.Decode { //If Decoding
		var p PngInfo
		offset, _ := strconv.ParseInt(md.Offset, 10, 64)
		buff.Seek(offset, 0)
		p.ChunkRead(buff)
		ogData := p.Chunk.Data
		p.Chunk.Data = baseutils.EncoderDecoder(p.Chunk.Data, md.Key)
		p.Chunk.CRC = p.CreateCRC()
		buffp := p.Marshal()
		buffpp := buffp.Bytes()
		fmt.Printf("Payload Original: % X\n", ogData)
		fmt.Printf("Payload Decode: % X\n", p.Chunk.Data)
		baseutils.WriteData(buff, md, buffpp)
	}
	if md.Meta {
		count := 1
		chunkType := ""
		endChunkType := "IEND" //last type represented before eof for png
		//Evaluates chunk type until eof for png
		for chunkType != endChunkType {
			fmt.Println("---- Chunk # " + strconv.Itoa(count) + " ----")
			pi.GetOffset(buff)
			fmt.Printf("Chunk Offset: %#02x\n", pi.Offset)
			pi.ChunkRead(buff)
			fmt.Printf("Chunk Length: %s bytes\n", strconv.Itoa(int(pi.Chunk.Size)))
			fmt.Printf("Chunk Type: %s\n", pi.chunkTypeToString())
			fmt.Printf("Chunk Importance: %s\n", pi.checkCritType())
			chunkType = pi.chunkTypeToString()
			count++
		}
	}
}
