package smfio

import (
	"bytes"
	"testing"

	"github.com/algoGuy/EasyMIDI/smf"
)

func TestReadNilReader(t *testing.T) {

	//act
	_, err := Read(nil)

	//assert
	notNilAssert(err, t)
}

func TestReadNoHeader(t *testing.T) {

	//arrange
	reader := &bytes.Buffer{}

	//act
	_, err := Read(reader)

	//assert
	notNilAssert(err, t)
}

func TestReadNoData(t *testing.T) {

	//arrange
	reader := bytes.NewBuffer(nil)

	//act
	_, err := Read(reader)

	//assert
	notNilAssert(err, t)
}

func TestReadCorruptedHeader(t *testing.T) {

	//arrange
	data := getMThd()
	data = append(data, []byte{0x00, 0x00, 0x00, 0x06}...)
	reader := bytes.NewBuffer(data)

	//act
	_, err := Read(reader)

	//assert
	notNilAssert(err, t)
}

func TestReadWrongFormat(t *testing.T) {

	//arrange
	data := getMThd()
	data = append(data, []byte{0x00, 0x00, 0x00, 0x06, 0x00, byte(smf.Format2 + 1), 0x00, 0x00, 0x00, 0x00}...)
	reader := bytes.NewBuffer(data)

	//act
	_, err := Read(reader)

	//assert
	notNilAssert(err, t)
}

func notNilAssert(value interface{}, t *testing.T) {

	if value == nil {
		t.Error("Wait not nil but was nil")
	}
}
