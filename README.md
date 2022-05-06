# finalproject
Code derived from Black Hat Go:Chapter 13

Usage: 
./finalproject
--User will be prompted for input throughout the program

Prompts:

FOR BOTH ENCODING AND DECODING

"Enter image location:"
Enter the location where the png image is located
Ex. apple.png

"Would you like to encode or decode this image(E/D)?:"

Enter E to encode and D to decode, can only do one at a time, defaults to encode if an invalid input is detected

Ex. D (capitalization does not matter) 

"Xor or Aes Method(X/A)?:"
Enter X for Xor encryption method and A for Aes encryption method, defaults to Xor method (Aes only works on strings of 16 bits as of now) 
Ex. X(capitalization does not matter) 

"Please enter your message key:"
Enter any BUT the empty string, will default to "key" if you do not enter a string
Ex. UWYO

"Pick an offset from the list above:" --NOT FUNCTIONAL 
Given a list of offsets, you can pick which one you would like to start at, defaults to the final critical chunk
Ex. 0x489a00

"Enter filename to save to:"
Defaults to encode.png when encoding and decode.png when decoding. 
Ex. appleEncode.png


FOR ENCODING

"Please enter your message to encode: "
If Xor, message can be any length, if Aes, currently must be 16 bits or longer. Will only encode first 16 bits however. 
Ex. I like apples!!!


FOR DECODING

"Please enter your offset:"
Enter the offset from encoding. 
Ex. 0x489a00

Example: 
Encoding

![image](https://user-images.githubusercontent.com/47127711/167201187-d551b45d-caf7-41c3-bffb-3ea179c11064.png)

Decoding

![image](https://user-images.githubusercontent.com/47127711/167201348-b2981166-b76c-446c-8c5b-8b65305d681f.png)


Current Successes
-Replaces bits instead of just dumping extra bits into the file 
-Partway implementation of Aes
-Scattering bits across chunks in png explored. Currently get incorrect offsets when using this option, exporing further. 
-JPEF and PDF files researches, further implementation needed. 
