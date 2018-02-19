package smfio

import (
	"encoding/binary"
	"io"

	"github.com/algoGuy/EasyMIDI/smf"
)

const minMThdDataSize uint32 = 6

//Getter for MThd header
func getMThd() []byte {
	return []byte("MThd")
}

//header MIDI File struct
type header struct {
	tracksNum uint16
	format    uint16
	division  *smf.Division
}

//uint32 representation of mthd
const mthdCode uint32 = 0x4D546864

//parseHeader parse header from reader
func parseHeader(reader io.Reader) (*header, error) {

	//parse mthdCode
	var mthd uint32
	if err := binary.Read(reader, binary.BigEndian, &mthd); err != nil {
		return nil, err
	}

	//check mthd
	if mthdCode != mthd {
		return nil, &ReadError{"No MThd"}
	}

	//get length
	var length uint32
	if err := binary.Read(reader, binary.BigEndian, &length); err != nil {
		return nil, err
	}

	//todo: add string format!
	if length < minMThdDataSize {
		return nil, &ReadError{"mthd length < 6"}
	}

	//read data
	data := make([]byte, length)
	if _, err := io.ReadAtLeast(reader, data, int(length)); err != nil {
		return nil, &ReadError{"Midi file corrupted"}
	}

	//parse header from data
	return parseHeaderData(data)
}

//parseHeader parse header from data
func parseHeaderData(data []byte) (*header, error) {

	//create midi
	format := binary.BigEndian.Uint16(data[0:2])
	tracksNum := binary.BigEndian.Uint16(data[2:4])

	//create division
	ticks := binary.BigEndian.Uint16(data[4:6])
	frames := int8(data[4])

	if frames > 0 {
		frames = smf.NOSMTPE
	} else {
		ticks &= 0xFF
	}

	//check supported division
	division, err := smf.NewDivision(ticks, frames)
	if err != nil {
		return nil, &ReadError{err.Error()}
	}

	return &header{format: format, tracksNum: tracksNum, division: division}, nil
}

//writeHeader writes header for midi to writer
func writeHeader(midi *smf.MIDIFile, writer io.Writer) {

	//add MThd string
	writer.Write(getMThd())

	//add size bytes
	writer.Write([]byte{0x00, 0x00, 0x00, byte(minMThdDataSize)})

	//add format
	writer.Write([]byte{0x00, uint8(midi.GetFormat())})

	//add tracks
	writer.Write(uint16toBytes(midi.GetTracksNum()))

	//add division
	division := midi.GetDivision()
	writer.Write(divisionToBytes(division))
}

//divisionToBytes creates []byte representation of division
func divisionToBytes(division smf.Division) []byte {

	//get bytes of ticks
	result := uint16toBytes(division.GetTicks())

	//if division
	if division.IsSMTPE() {
		result[0] = byte(division.GetSMTPE())
	}

	return result
}

//uint16toBytes translate uint16 to []byte
func uint16toBytes(value uint16) []byte {

	//create result
	result := make([]byte, 2)
	binary.BigEndian.PutUint16(result, value)

	return result
}
