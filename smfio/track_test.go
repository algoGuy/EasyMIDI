package smfio

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/algoGuy/EasyMIDI/smf"
	"github.com/algoGuy/EasyMIDI/vlq"
)

func TestIsEndOfTrack(t *testing.T) {

	//arrange
	event, _ := smf.NewMetaEvent(0, smf.MetaEndOfTrack, []byte{0x00})

	//act
	result := isEndOfTrack(event)

	//assert
	if !result {
		t.Error("wait true but was false")
	}
}

func TestIsEndOfTrackWrong(t *testing.T) {

	//arrange
	firstEvent, _ := smf.NewMetaEvent(0, smf.MetaCuePoint, []byte{0x00})
	secondEvent, _ := smf.NewMIDIEvent(0, smf.NoteOnStatus, 0, 0x00, 0x00)
	thirdEvent, _ := smf.NewSysexEvent(0, smf.SysexStatus, []byte{0x00})

	testData := []smf.Event{
		nil,
		firstEvent,
		secondEvent,
		thirdEvent,
	}

	for _, event := range testData {

		//act
		result := isEndOfTrack(event)

		//assert
		if result {
			t.Error("wait false but was true")
		}
	}
}

func TestWritTrack(t *testing.T) {

	//arrange
	track := &smf.Track{}

	midiEventTime := vlq.MaxNumber / 2
	midiEvent, _ := smf.NewMIDIEvent(midiEventTime, smf.NoteOffStatus, 0, 0x00, 0x00)
	track.AddEvent(midiEvent)

	endEvenTime := vlq.MaxNumber
	endEvent, _ := smf.NewMetaEvent(endEvenTime, smf.MetaEndOfTrack, []byte{})
	track.AddEvent(endEvent)

	writer := &bytes.Buffer{}

	//act
	err := writeTrack(track, writer)

	//assert
	if err != nil {
		t.Error("wait no error but was error")
	}

	realLength := uint32(6 + len(vlq.GetBytes(midiEventTime)) + len(vlq.GetBytes(endEvenTime)))

	realResult := getMTrk()
	realResult = append(realResult, uint32toBytes(realLength)...)

	realResult = append(realResult, vlq.GetBytes(midiEventTime)...)
	realResult = append(realResult, midiEvent.GetChannel()|midiEvent.GetStatus())
	realResult = append(realResult, midiEvent.GetData()...)

	realResult = append(realResult, vlq.GetBytes(endEvenTime)...)
	realResult = append(realResult, endEvent.GetStatus())
	realResult = append(realResult, endEvent.GetMetaType())
	realResult = append(realResult, vlq.GetBytes(uint32(len(endEvent.GetData())))...)
	realResult = append(realResult, endEvent.GetData()...)

	if !bytes.Equal(realResult, writer.Bytes()) {
		t.Errorf("Wait % x but was % x", realResult, writer.Bytes())
	}
}

func TestWritTrackNilTrack(t *testing.T) {

	//act
	err := writeTrack(nil, nil)

	//assert
	if err == nil {
		t.Error("wait err but was nil")
	}
}

func TestParseTrack(t *testing.T) {

	//arrange
	//add sysex
	sysexEvent, _ := smf.NewSysexEvent(vlq.MaxNumber, smf.SysexStatus, []byte{0x00, 0x01})
	data := vlq.GetBytes(sysexEvent.GetDTime())
	data = append(data, sysexEvent.GetStatus())
	data = append(data, vlq.GetBytes(uint32(len(sysexEvent.GetData())))...)
	data = append(data, sysexEvent.GetData()...)

	//add midi event
	midiEvent, _ := smf.NewMIDIEvent(0, smf.NoteOnStatus, smf.MaxChannelNumber, 0x00, 0x01)
	data = append(data, vlq.GetBytes(midiEvent.GetDTime())...)
	data = append(data, midiEvent.GetStatus()|midiEvent.GetChannel())
	data = append(data, midiEvent.GetData()...)

	//add meta event
	metaEvent, _ := smf.NewMetaEvent(vlq.MaxNumber-1, smf.MetaEndOfTrack, []byte{})
	data = append(data, vlq.GetBytes(metaEvent.GetDTime())...)
	data = append(data, metaEvent.GetStatus())
	data = append(data, metaEvent.GetMetaType())
	data = append(data, vlq.GetBytes(uint32(len(metaEvent.GetData())))...)
	data = append(data, midiEvent.GetData()...)

	reader := bytes.NewReader(data)

	realTrack := &smf.Track{}
	realTrack.AddEvent(sysexEvent)
	realTrack.AddEvent(midiEvent)
	realTrack.AddEvent(metaEvent)

	//act
	track, err := parseTrack(reader)

	//assert
	if err != nil {
		t.Error("wait nil but was error")
	}

	if !reflect.DeepEqual(realTrack, track) {
		t.Errorf("wait %v but was %v", realTrack, track)
	}
}

func TestParseTrackNoEnd(t *testing.T) {

	//arrage
	//add midi event
	midiEvent, _ := smf.NewMIDIEvent(0, smf.NoteOnStatus, smf.MaxChannelNumber, 0x00, 0x01)
	data := vlq.GetBytes(midiEvent.GetDTime())
	data = append(data, midiEvent.GetStatus()|midiEvent.GetChannel())
	data = append(data, midiEvent.GetData()...)

	reader := bytes.NewReader(data)

	//act
	_, err := parseTrack(reader)

	//assert
	if err == nil {
		t.Error("wait err but was nil")
	}
}

func TestParseAllTracks(t *testing.T) {

	//add midi event
	midiEvent, _ := smf.NewMIDIEvent(0, smf.NoteOnStatus, smf.MaxChannelNumber, 0x00, 0x01)
	data := vlq.GetBytes(midiEvent.GetDTime())
	data = append(data, midiEvent.GetStatus()|midiEvent.GetChannel())
	data = append(data, midiEvent.GetData()...)

	//add meta event
	metaEvent, _ := smf.NewMetaEvent(vlq.MaxNumber-1, smf.MetaEndOfTrack, []byte{})
	data = append(data, vlq.GetBytes(metaEvent.GetDTime())...)
	data = append(data, metaEvent.GetStatus())
	data = append(data, metaEvent.GetMetaType())
	data = append(data, vlq.GetBytes(uint32(len(metaEvent.GetData())))...)
	data = append(data, metaEvent.GetData()...)

	//add prefix
	data = append(uint32toBytes(uint32(len(data))), data...)
	data = append(getMTrk(), data...)

	//double track
	data = append(data, data...)

	division, _ := smf.NewDivision(120, 0)
	testMidi, _ := smf.NewSMF(smf.Format1, *division)

	//act
	midi, err := parseAllTracks(testMidi, bytes.NewReader(data), 2)

	//assert
	realMidi, _ := smf.NewSMF(smf.Format1, *division)

	track := &smf.Track{}
	track.AddEvent(midiEvent)
	track.AddEvent(metaEvent)

	realMidi.AddTrack(track)
	realMidi.AddTrack(track)

	if err != nil {
		t.Error("wait nil but was error")
	}

	if !reflect.DeepEqual(midi, realMidi) {
		t.Errorf("wait %v but was %v", midi, realMidi)
	}
}

func TestParseAllTracksNoData(t *testing.T) {

	//arrage
	reader := bytes.NewReader(nil)
	midi := &smf.MIDIFile{}

	//act
	_, err := parseAllTracks(midi, reader, 1)

	//assert
	if err == nil {
		t.Error("wait error but was nil")
	}
}

func TestParseAllTracksNoHeader(t *testing.T) {

	//arrange
	data := []byte{0x00, 0x02, 0x04, 0x06}
	midi := &smf.MIDIFile{}

	//act
	_, err := parseAllTracks(midi, bytes.NewReader(data), 1)

	//assert
	if err == nil {
		t.Error("wait error but was nil")
	}
}

func TestParseAllTracksNoLength(t *testing.T) {

	//arrange
	data := append(getMTrk(), 0x00)
	midi := &smf.MIDIFile{}

	//act
	_, err := parseAllTracks(midi, bytes.NewReader(data), 1)

	//assert
	if err == nil {
		t.Error("wait error but was nil")
	}
}

func TestParseAllTracksTrackError(t *testing.T) {

	//arrange
	data := append(getMTrk(), []byte{0x00, 0x00, 0x00, 0x02, 0x00}...)
	midi := &smf.MIDIFile{}

	//act
	_, err := parseAllTracks(midi, bytes.NewReader(data), 1)

	//assert
	if err == nil {
		t.Error("wait error but was nil")
	}
}

func TestParseAllTracksAddTrackError(t *testing.T) {

	//arrange
	//add meta event
	metaEvent, _ := smf.NewMetaEvent(vlq.MaxNumber-1, smf.MetaEndOfTrack, []byte{})
	data := vlq.GetBytes(metaEvent.GetDTime())
	data = append(data, metaEvent.GetStatus())
	data = append(data, metaEvent.GetMetaType())
	data = append(data, vlq.GetBytes(uint32(len(metaEvent.GetData())))...)
	data = append(data, metaEvent.GetData()...)

	//add prefix
	data = append(uint32toBytes(uint32(len(data))), data...)
	data = append(getMTrk(), data...)

	//double track
	data = append(data, data...)

	division, _ := smf.NewDivision(120, 0)
	testMidi, _ := smf.NewSMF(smf.Format0, *division)

	//act
	_, err := parseAllTracks(testMidi, bytes.NewReader(data), 2)

	//assert
	if err == nil {
		t.Error("wait error but was nil")
	}
}

func TestWriteAllTracks(t *testing.T) {

	//arrange
	//add midi event
	midiEvent, _ := smf.NewMIDIEvent(0, smf.NoteOnStatus, smf.MaxChannelNumber, 0x00, 0x01)
	data := vlq.GetBytes(midiEvent.GetDTime())
	data = append(data, midiEvent.GetStatus()|midiEvent.GetChannel())
	data = append(data, midiEvent.GetData()...)

	//add meta event
	metaEvent, _ := smf.NewMetaEvent(vlq.MaxNumber-1, smf.MetaEndOfTrack, []byte{})
	data = append(data, vlq.GetBytes(metaEvent.GetDTime())...)
	data = append(data, metaEvent.GetStatus())
	data = append(data, metaEvent.GetMetaType())
	data = append(data, vlq.GetBytes(uint32(len(metaEvent.GetData())))...)
	data = append(data, metaEvent.GetData()...)

	//add prefix
	data = append(uint32toBytes(uint32(len(data))), data...)
	data = append(getMTrk(), data...)

	//double data
	data = append(data, data...)

	division, _ := smf.NewDivision(120, 0)
	realMidi, _ := smf.NewSMF(smf.Format1, *division)

	track := &smf.Track{}
	track.AddEvent(midiEvent)
	track.AddEvent(metaEvent)

	realMidi.AddTrack(track)
	realMidi.AddTrack(track)

	writer := &bytes.Buffer{}

	//act
	err := writeAllTracks(realMidi, writer)

	//assert
	if err != nil {
		t.Error("wait for nil but was error")
	}

	if !bytes.Equal(data, writer.Bytes()) {
		t.Errorf("wait for % x but was % x", data, writer.Bytes())
	}
}

func TestUint32ToBytes(t *testing.T) {

	//arrange
	testData := []uint32{0, 1, 1 << 8, 1 << 16, 1 << 24, 1 << 31}

	//tests
	for _, data := range testData {

		//act
		result := uint32toBytes(data)

		//assert
		realResult := make([]byte, 4)
		binary.BigEndian.PutUint32(realResult, 4)

		if bytes.Equal(realResult, result) {
			t.Errorf("Wait % x but was % x", realResult, result)
		}
	}
}
