package baseutils

//XOR D/Encode payload from key string
func EncoderDecoder(in []byte, key string) []byte {
	buffArr := make([]byte, len(in))
	for i := 0; i < len(in); i++ {
		buffArr[i] += in[i] ^ key[i%len(key)]
	}
	return buffArr
}
