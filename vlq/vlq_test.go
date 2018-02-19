package vlq

import (
	"bytes"
	"math"
	"reflect"
	"testing"
)

func TestGetBytesFromVLQ(t *testing.T) {

	//arrange
	vlqTests := []struct {
		in  uint32
		out []byte
	}{
		{0x00, []byte{0x00}},
		{0x01, []byte{0x01}},
		{0xFF, []byte{0x81, 0x7F}},
		{0xFFFF, []byte{0x83, 0xFF, 0x7F}},
		{0xFFFFFF, []byte{0x87, 0xFF, 0xFF, 0x7F}},
		{MaxNumber, []byte{0xFF, 0xFF, 0xFF, 0x7F}},
		{uint32(math.MaxUint32), nil},
		{MaxNumber + 1, nil},
	}

	for _, testData := range vlqTests {

		//act
		result := GetBytes(testData.in)

		//assert
		if !reflect.DeepEqual(result, testData.out) {
			t.Errorf("Test case %X wait % X but was % X", testData.in, testData.out, result)
		}
	}
}

//arrange
var vlqTests = []struct {
	in  []byte
	out uint32
}{
	{[]byte{0x00}, 0x00},
	{[]byte{0x01}, 0x01},
	{[]byte{0x81, 0x7F}, 0xFF},
	{[]byte{0x83, 0xFF, 0x7F}, 0xFFFF},
	{[]byte{0x87, 0xFF, 0xFF, 0x7F}, 0xFFFFFF},
	{[]byte{0xFF, 0xFF, 0xFF, 0x7F}, MaxNumber},
}

func TestGetFromBytes(t *testing.T) {

	for _, testData := range vlqTests {

		//act
		result, length := GetFromBytes(testData.in)

		//assert
		if result != testData.out {
			t.Errorf("Test case % X wait % X but was % X", testData.in, testData.out, result)
		}

		if int(length) != len(testData.in) {
			t.Errorf("Test case % X wait length %X but was %X", testData.in, len(testData.in), result)
		}
	}
}

func TestGetFromReader(t *testing.T) {

	for _, testData := range vlqTests {

		//arrange
		reader := bytes.NewReader(testData.in)

		//act
		result, err := GetFromReader(reader)

		//assert
		if result != testData.out {
			t.Errorf("Test case % X wait %X but was %X", testData.in, testData.out, result)
		}

		if err != nil {
			t.Error("Unexpected error")
		}
	}
}

func TestGetFromReaderMsb1(t *testing.T) {

	//arrange
	testData := [][]byte{
		[]byte{Msb1Mask},
		[]byte{Msb1Mask, Msb1Mask},
		[]byte{Msb1Mask, Msb1Mask, Msb1Mask},
		[]byte{Msb1Mask, Msb1Mask, Msb1Mask, Msb1Mask},
	}

	//tests
	for _, data := range testData {

		reader := bytes.NewBuffer(data)
		if _, err := GetFromReader(reader); err == nil {
			t.Errorf("wait error but was nil on % x", data)
		}
	}
}
