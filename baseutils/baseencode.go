package baseutils

import (
	"crypto/aes"
	"log"
)

//XOR D/Encode payload from key string
func XorEncoderDecoder(in []byte, key string) []byte {
	buffArr := make([]byte, len(in))
	for i := 0; i < len(in); i++ {
		buffArr[i] += in[i] ^ key[i%len(key)]
	}
	return buffArr
}

//Aes D/Encode payload from key string --NEEDS WORK
func AesEncoderDecoder(in []byte, key string, code string) []byte {
	space := "a"
	buffArr := make([]byte, len(in))
	if len(key) > 32 {
		key = key[0:31]
	}
	for len(key) != 32 {
		key = key + space
	}

	ciph, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Fatal(err)
	}

	if code == "E" {
		ciph.Encrypt(buffArr, in)
	} else {
		ciph.Decrypt(buffArr, in)
	}
	return buffArr
}

func EncoderDecoder(in []byte, key string, ty string, code string) []byte {
	if ty == "X" {
		return XorEncoderDecoder(in, key)
	} else {
		return AesEncoderDecoder(in, key, code)
	}
}
