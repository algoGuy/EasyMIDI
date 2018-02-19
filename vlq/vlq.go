package vlq

import (
	"encoding/binary"
	"io"
)

//MaxNumber - maximum VLQ number
const MaxNumber uint32 = 0xFFFFFFF

//vlqSegmentSize bit size of vlq number
const vlqSegmentSize = 7

//MaxResultSize = maximum bytes number for vlq
const MaxResultSize = 4

//bits masks for vlq
const (
	Msb1Mask   byte = 0x80
	MsbNo1Mask byte = 0x7F
)

//GetFromBytes gets Variable-length quantity from data
func GetFromBytes(data []byte) (deltaTime uint32, size uint32) {

	for ; size < uint32(len(data)); size++ {

		//add next byte
		deltaTime <<= vlqSegmentSize
		deltaTime |= uint32(data[size] & MsbNo1Mask)

		//check 1 msb
		if data[size]&Msb1Mask == 0 {
			break
		}
	}

	return deltaTime, size + 1
}

//GetBytes return bytes representation of number
func GetBytes(number uint32) []byte {

	//check max size
	if number > MaxNumber {
		return nil
	}

	//create result
	const MaxResultSize = 4
	var result [MaxResultSize]byte

	//construct loop WHERE DO WHILE?! TODO: ADD TO GO DO WHILE LOOPS
	start := MaxResultSize - 1
	for ; ; start-- {

		result[start] = byte(number&0X7F) | Msb1Mask
		number >>= vlqSegmentSize

		if number == 0 {
			break
		}
	}

	//set first msb to 0
	result[MaxResultSize-1] &= MsbNo1Mask
	return result[start:]
}

//GetFromReader get vlq from reader
func GetFromReader(reader io.Reader) (uint32, error) {

	var deltaTime uint32
	var nextByte byte
	for i := 0; i < MaxResultSize; i++ {

		//get next byte
		if err := binary.Read(reader, binary.BigEndian, &nextByte); err != nil {
			return 0, err
		}

		//calculate dtime
		deltaTime <<= vlqSegmentSize
		deltaTime |= uint32(nextByte & MsbNo1Mask)

		//check nextByte msb == 0
		if nextByte < Msb1Mask {
			break
		}
	}

	//check last byte
	if nextByte >= Msb1Mask {
		return 0, &Error{"last byte has msb == 1"}
	}

	return deltaTime, nil
}
