package smfio

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/algoGuy/EasyMIDI/smf"
)

//Getter for MTrk id
func getMTrk() []byte {
	return []byte("MTrk")
}

const mtrkCode uint32 = 0x4D54726B

//parseAllTracks parse all tracks from reader
func parseAllTracks(midi *smf.MIDIFile, reader io.Reader, tracksNum uint16) (*smf.MIDIFile, error) {

	for i := uint16(0); i < tracksNum; i++ {

		//parse mtrkCode
		var mtrk uint32
		if err := binary.Read(reader, binary.BigEndian, &mtrk); err != nil {
			return nil, &ReadError{"Track data corrupted"}
		}

		//check mthd
		if mtrkCode != mtrk {
			return nil, &ReadError{"No MTrk on track start"}
		}

		//get length
		var length uint32
		if err := binary.Read(reader, binary.BigEndian, &length); err != nil {
			return nil, &ReadError{"Track data corrupted"}
		}

		//create limited reader
		trackReader := &io.LimitedReader{R: reader, N: int64(length)}

		//parse track
		track, err := parseTrack(trackReader)
		if err != nil {
			return nil, err
		}

		//add track + check error
		err = midi.AddTrack(track)
		if err != nil {
			return nil, err
		}
	}

	return midi, nil
}

//parseTrack parse one track with events
func parseTrack(reader io.Reader) (*smf.Track, error) {

	//create track
	track := &smf.Track{}

	//parse events
	var event smf.Event
	var err error

	for !isEndOfTrack(event) {

		event, err = parseEvent(reader, event)
		if err != nil {
			return nil, &ReadError{err.Error()}
		}

		track.AddEvent(event)
	}

	return track, nil
}

//writeAllTracks writes midi tracks to writer
func writeAllTracks(midi *smf.MIDIFile, writer io.Writer) error {

	//start analyze loop
	tracksNum := midi.GetTracksNum()

	for i := uint16(0); i < tracksNum; i++ {

		//start write track
		err := writeTrack(midi.GetTrack(i), writer)

		//check error
		if err != nil {
			return err
		}
	}

	return nil
}

//writeTrack write track to
func writeTrack(track *smf.Track, writer io.Writer) error {

	//check nil reference
	if track == nil {
		return &WriteError{"MIDI corrupted"}
	}

	//add track chunk name
	writer.Write(getMTrk())

	//create binary writer
	buffer := &bytes.Buffer{}

	//start iterate
	var prevEvent smf.Event
	iterator := track.GetIterator()
	for iterator.MoveNext() {

		//write events
		event := iterator.GetValue()
		err := writeEvent(event, prevEvent, buffer)

		//check error
		if err != nil {
			return err
		}

		//set prev event
		prevEvent = event
	}

	//write len
	writer.Write(uint32toBytes(uint32(buffer.Len())))
	writer.Write(buffer.Bytes())

	return nil
}

//isEndOfTrack is corrent event EOT
func isEndOfTrack(event smf.Event) bool {

	switch v := event.(type) {
	case *smf.MetaEvent:
		return v.GetMetaType() == smf.MetaEndOfTrack
	default:
		return false
	}
}

//uint32toBytes translate uint16 to []byte
func uint32toBytes(value uint32) []byte {

	//create result
	result := make([]byte, 4)
	binary.BigEndian.PutUint32(result, value)

	return result
}
