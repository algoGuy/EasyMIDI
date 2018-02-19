package smfio

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/algoGuy/EasyMIDI/smf"
)

func TestDivisionToBytes(t *testing.T) {

	//arrange
	ticks := smf.TicksMaxValue
	SMTPEFrames := smf.NOSMTPE

	division, _ := smf.NewDivision(ticks, SMTPEFrames)

	//act
	result := divisionToBytes(*division)

	//assert
	realResult := []byte{byte(ticks >> 8), byte(ticks)}
	if !bytes.Equal(result, realResult) {
		t.Errorf("First array %v not equal with %v", result, realResult)
	}
}

func TestDivisionToBytesSMTPE(t *testing.T) {

	//arrange
	ticks := smf.SMTPETicksMaxValue
	SMTPEFrames := smf.SMTPE29

	division, _ := smf.NewDivision(ticks, SMTPEFrames)

	//act
	result := divisionToBytes(*division)

	//assert
	realResult := []byte{byte(SMTPEFrames), byte(ticks)}
	if !bytes.Equal(result, realResult) {
		t.Errorf("First array %v not equal with %v", result, realResult)
	}
}

func TestWriteHeader(t *testing.T) {

	//arrange
	format := smf.Format1
	ticks := smf.SMTPETicksMaxValue
	SMTPEFrames := smf.SMTPE29
	tracksNum := uint16(10)

	division, _ := smf.NewDivision(ticks, SMTPEFrames)
	midiFile, _ := smf.NewSMF(format, *division)

	for i := uint16(0); i < tracksNum; i++ {
		midiFile.AddTrack(&smf.Track{})
	}

	writer := &bytes.Buffer{}

	//act
	writeHeader(midiFile, writer)

	//assert
	realResult := getMThd()
	realResult = append(realResult, []byte{0x00, 0x00, 0x00, byte(minMThdDataSize), 0x00, byte(format)}...)
	realResult = append(realResult, uint16toBytes(tracksNum)...)
	realResult = append(realResult, divisionToBytes(*division)...)

	if !bytes.Equal(writer.Bytes(), realResult) {
		t.Errorf("First array %v not equal with %v", writer.Bytes(), realResult)
	}
}

func TestParseHeaderData(t *testing.T) {

	//arrange
	format := smf.Format1
	ticks := smf.TicksMaxValue
	SMTPEFrames := smf.NOSMTPE
	tracksNum := uint16(10)

	division, _ := smf.NewDivision(ticks, SMTPEFrames)

	//create test data
	data := []byte{0x00, byte(format)}
	data = append(data, uint16toBytes(tracksNum)...)
	data = append(data, divisionToBytes(*division)...)

	//act
	result, err := parseHeaderData(data)

	//assert
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(result.division, division) {
		t.Error("Wrong division parsing")
	}

	if result.format != format {
		t.Error("Wrong format parsing")
	}

	if result.tracksNum != tracksNum {
		t.Error("Wrong trackNum parsing")
	}
}

func TestParseHeaderDataSMTPE(t *testing.T) {

	//arrange
	format := smf.Format1
	ticks := smf.SMTPETicksMaxValue
	SMTPEFrames := smf.SMTPE24
	tracksNum := uint16(10)

	division, _ := smf.NewDivision(ticks, SMTPEFrames)

	//create test data
	data := []byte{0x00, byte(format)}
	data = append(data, uint16toBytes(tracksNum)...)
	data = append(data, divisionToBytes(*division)...)

	//act
	result, err := parseHeaderData(data)

	//assert
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(result.division, division) {
		t.Error("Wrong division parsing")
	}

	if result.format != format {
		t.Error("Wrong format parsing")
	}

	if result.tracksNum != tracksNum {
		t.Error("Wrong trackNum parsing")
	}
}

func TestParseHeaderDataWrongTicks(t *testing.T) {

	//arrange
	format := smf.Format1
	ticks := smf.TicksMaxValue + 1
	tracksNum := uint16(10)

	//create test data
	data := []byte{0x00, byte(format)}
	data = append(data, uint16toBytes(tracksNum)...)
	data = append(data, uint16toBytes(ticks)...)

	//act
	_, err := parseHeaderData(data)

	//assert
	if err == nil {
		t.Error("Wait error but was nil")
	}
}

func TestParseHeaderNoMThd(t *testing.T) {

	//arrange
	data := []byte{0x00, 0x00, 0x00, 0x00}
	reader := bytes.NewBuffer(data)

	//act
	_, err := parseHeader(reader)

	//assert
	if err == nil {
		t.Error("wait error but was nil")
	}
}

func TestParseHeaderNoLength(t *testing.T) {

	//arrange
	data := getMThd()
	reader := bytes.NewBuffer(data)

	//act
	_, err := parseHeader(reader)

	//assert
	if err == nil {
		t.Error("wait error but was nil")
	}
}

func TestParseHeaderWrongLength(t *testing.T) {

	//arrange
	data := getMThd()
	data = append(data, []byte{0x00, 0x00, 0x00, 0x00}...)
	reader := bytes.NewBuffer(data)

	//act
	_, err := parseHeader(reader)

	//assert
	if err == nil {
		t.Error("wait error but was nil")
	}
}
