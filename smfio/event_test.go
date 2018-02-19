package smfio

import (
	"bytes"
	"reflect"
	"testing"

	"bitbucket.org/NewStreeter/MIDIParser/smf"
	"bitbucket.org/NewStreeter/MIDIParser/vlq"
)

func TestWriteEventMeta(t *testing.T) {

	//arrange
	deltaTime := vlq.MaxNumber
	metaType := smf.MetaEndOfTrack
	data := []byte{0x00, 0x01}
	metaEvent, _ := smf.NewMetaEvent(deltaTime, metaType, data)

	writer := &bytes.Buffer{}

	//act
	writeEvent(metaEvent, nil, writer)

	//assert
	realResult := vlq.GetBytes(deltaTime)
	realResult = append(realResult, smf.MetaStatus)
	realResult = append(realResult, metaType)
	realResult = append(realResult, vlq.GetBytes(uint32(len(data)))...)
	realResult = append(realResult, data...)

	if !bytes.Equal(realResult, writer.Bytes()) {
		t.Errorf("Wait for % x array but was % x array", realResult, writer.Bytes())
	}
}

func TestWriteEventSysex(t *testing.T) {

	//arrange
	deltaTime := vlq.MaxNumber

	status := smf.SysexStatus
	data := []byte{0x00, 0x01}
	metaEvent, _ := smf.NewSysexEvent(deltaTime, status, data)

	writer := &bytes.Buffer{}

	//act
	writeEvent(metaEvent, nil, writer)

	//assert
	realResult := vlq.GetBytes(deltaTime)
	realResult = append(realResult, status)
	realResult = append(realResult, vlq.GetBytes(uint32(len(data)))...)
	realResult = append(realResult, data...)

	if !bytes.Equal(realResult, writer.Bytes()) {
		t.Errorf("Wait for % x array but was % x array", realResult, writer.Bytes())
	}
}

func TestWriteEventMidi(t *testing.T) {

	//arrange
	deltaTime := vlq.MaxNumber

	status := smf.NoteOnStatus
	channel := byte(0x01)
	event, _ := smf.NewMIDIEvent(deltaTime, status, channel, 0x00, 0x01)

	writer := &bytes.Buffer{}

	//act
	writeEvent(event, nil, writer)

	//assert
	MidiEventAssert(event, writer, false, t)
}

func TestWriteEventRunningStatus(t *testing.T) {

	//arrange
	deltaTime := vlq.MaxNumber

	status := smf.NoteOnStatus
	channel := byte(0x01)
	event, _ := smf.NewMIDIEvent(deltaTime, status, channel, 0x00, 0x01)

	writer := &bytes.Buffer{}

	//act
	writeEvent(event, event, writer)

	//assert
	MidiEventAssert(event, writer, true, t)
}

func TestWriteEventRunningStatusWrongChannel(t *testing.T) {

	//arrange
	deltaTime := vlq.MaxNumber

	status := smf.NoteOnStatus
	channel := byte(0x01)
	event, _ := smf.NewMIDIEvent(deltaTime, status, channel, 0x00, 0x01)
	event2, _ := smf.NewMIDIEvent(deltaTime, status, channel+1, 0x00, 0x01)

	writer := &bytes.Buffer{}

	//act
	writeEvent(event, event2, writer)

	//assert
	MidiEventAssert(event, writer, false, t)
}

func TestWriteEventRunningStatusWrongStatus(t *testing.T) {

	//arrange
	deltaTime := vlq.MaxNumber

	status := smf.NoteOnStatus
	status2 := smf.NoteOffStatus

	channel := byte(0x01)
	event, _ := smf.NewMIDIEvent(deltaTime, status, channel, 0x00, 0x01)
	event2, _ := smf.NewMIDIEvent(deltaTime, status2, channel, 0x00, 0x01)

	writer := &bytes.Buffer{}

	//act
	writeEvent(event, event2, writer)

	//assert
	MidiEventAssert(event, writer, false, t)
}

func MidiEventAssert(event *smf.MIDIEvent, writer *bytes.Buffer, runningStatus bool, t *testing.T) {

	//calculate real result
	realResult := vlq.GetBytes(event.GetDTime())

	if !runningStatus {
		realResult = append(realResult, event.GetStatus()|event.GetChannel())
	}

	realResult = append(realResult, event.GetData()...)

	//assert
	if !bytes.Equal(realResult, writer.Bytes()) {
		t.Errorf("Wait for % x array but was % x array", realResult, writer.Bytes())
	}
}

func TestParseEventMidi(t *testing.T) {

	//arrange
	status := smf.NoteOnStatus
	channel := smf.MaxChannelNumber
	dataBytes := []byte{0x00, 0x01}

	dTime := vlq.MaxNumber

	data := vlq.GetBytes(dTime)
	data = append(data, status|channel)
	data = append(data, dataBytes...)

	//act-assert
	realEvent, _ := smf.NewMIDIEvent(dTime, status, channel, dataBytes[0], dataBytes[1])
	EventAssert(data, realEvent, nil, t)
}

func TestParseEventMidiOneByte(t *testing.T) {

	//arrange
	status := smf.ProgramChangeStatus
	channel := smf.MaxChannelNumber
	dataBytes := []byte{0x02, 0x00}

	dTime := vlq.MaxNumber

	data := vlq.GetBytes(dTime)
	data = append(data, status|channel)
	data = append(data, dataBytes...)

	realEvent, _ := smf.NewMIDIEvent(dTime, status, channel, dataBytes[0], 0x00)

	//act-assert
	EventAssert(data, realEvent, nil, t)
}

func TestParseEventMidiOutOneOfRange(t *testing.T) {

	//arrange
	data := vlq.GetBytes(vlq.MaxNumber)
	data = append(data, smf.NoteOnStatus)

	//act-assert
	EventAssert(data, nil, nil, t)
}

func TestParseEventMidiOutTwoOfRange(t *testing.T) {

	//arrange
	data := vlq.GetBytes(vlq.MaxNumber)
	data = append(data, []byte{smf.NoteOnStatus, 0x00}...)

	//act-assert
	EventAssert(data, nil, nil, t)
}

func TestParseEventSysex(t *testing.T) {

	//arrange
	dTime := uint32(0)
	status := smf.SysexStatus
	eventData := []byte{0x01, 0x01, 0x01, 0x02, 0x03}
	eventDataLen := uint32(len(eventData))

	data := vlq.GetBytes(dTime)
	data = append(data, status)
	data = append(data, vlq.GetBytes(eventDataLen)...)
	data = append(data, eventData...)

	realEvent, _ := smf.NewSysexEvent(dTime, status, eventData)

	//act - assert
	EventAssert(data, realEvent, nil, t)
}

func TestParseEventSysexWrongLength(t *testing.T) {

	//arrange
	testCase1 := vlq.GetBytes(0)
	testCase1 = append(testCase1, smf.SysexStatus)

	testCase2 := append(testCase1, 0x01)
	testCase3 := append(testCase1, []byte{0x02, 0x00}...)

	testData := [][]byte{
		testCase1,
		testCase2,
		testCase3,
	}

	//tests
	for _, data := range testData {

		//act-assert
		EventAssert(data, nil, nil, t)
	}
}

func TestParseEventMeta(t *testing.T) {

	//arrange
	dTime := vlq.MaxNumber
	metaType := smf.MetaEndOfTrack
	eventData := []byte{0x00}
	eventDataLen := uint32(len(eventData))

	data := vlq.GetBytes(dTime)
	data = append(data, []byte{smf.MetaStatus, metaType}...)
	data = append(data, vlq.GetBytes(eventDataLen)...)
	data = append(data, eventData...)

	//act-assert
	realEvent, _ := smf.NewMetaEvent(dTime, metaType, eventData)
	EventAssert(data, realEvent, nil, t)
}

func TestParseEventMetaWrongType(t *testing.T) {

	//arrange
	data := vlq.GetBytes(vlq.MaxNumber)
	data = append(data, []byte{smf.MetaStatus, smf.MetaSequencerSpecific + 1, 0x00}...)

	//act-assert
	EventAssert(data, nil, nil, t)
}

func TestParseEventMetaWrongLength(t *testing.T) {

	//arrange
	testCase1 := vlq.GetBytes(vlq.MaxNumber)
	testCase1 = append(testCase1, smf.MetaStatus)

	testCase2 := append(testCase1, smf.MetaCuePoint)
	testCase3 := append(testCase2, []byte{0x02, 0x00}...)

	testData := [][]byte{
		testCase1,
		testCase2,
		testCase3,
	}

	//tests
	for _, data := range testData {

		//act-assert
		EventAssert(data, nil, nil, t)
	}
}

func TestParseEventWrongStatus(t *testing.T) {

	//arrange
	data := vlq.GetBytes(vlq.MaxNumber)
	data = append(data, 0x00)

	//act - assert
	EventAssert(data, nil, nil, t)
}

func TestParseEventRunningStatus(t *testing.T) {

	//arrange
	status := smf.NoteOnStatus
	channel := smf.MaxChannelNumber
	dataBytes := []byte{0x02, 0x00}

	dTime := vlq.MaxNumber
	data := vlq.GetBytes(dTime)
	data = append(data, dataBytes...)

	realEvent, _ := smf.NewMIDIEvent(dTime, status, channel, dataBytes[0], dataBytes[1])

	//act - assert
	EventAssert(data, realEvent, realEvent, t)
}

func TestParseEventRunningOneByte(t *testing.T) {

	//arrange
	status := smf.ProgramChangeStatus
	channel := smf.MaxChannelNumber
	dataBytes := []byte{0x02}

	dTime := vlq.MaxNumber
	data := vlq.GetBytes(dTime)
	data = append(data, dataBytes...)

	realEvent, _ := smf.NewMIDIEvent(dTime, status, channel, dataBytes[0], 0x00)

	//act - assert
	EventAssert(data, realEvent, realEvent, t)
}

func TestParseEventRunningNextOneMSB(t *testing.T) {

	//arrange
	status := smf.ProgramChangeStatus
	channel := smf.MaxChannelNumber
	dataBytes := []byte{vlq.Msb1Mask}

	dTime := vlq.MaxNumber
	data := vlq.GetBytes(dTime)
	data = append(data, dataBytes...)

	prevEvent, _ := smf.NewMIDIEvent(dTime, status, channel, 0x00, 0x00)

	//act - assert
	EventAssert(data, nil, prevEvent, t)
}

func TestParseEventRunningNoSecondByte(t *testing.T) {

	//arrange
	status := smf.NoteOffStatus
	channel := smf.MaxChannelNumber
	dataBytes := []byte{vlq.MsbNo1Mask}

	dTime := vlq.MaxNumber
	data := vlq.GetBytes(dTime)
	data = append(data, dataBytes...)

	prevEvent, _ := smf.NewMIDIEvent(dTime, status, channel, 0x00, 0x00)

	//act - assert
	EventAssert(data, nil, prevEvent, t)
}

func TestParseEventRunningStatusNextSecondMSB(t *testing.T) {

	//arrange
	status := smf.NoteOffStatus
	channel := smf.MaxChannelNumber
	dataBytes := []byte{0x00, vlq.Msb1Mask}

	dTime := vlq.MaxNumber
	data := vlq.GetBytes(dTime)
	data = append(data, dataBytes...)

	prevEvent, _ := smf.NewMIDIEvent(dTime, status, channel, 0x00, 0x00)

	//act - assert
	EventAssert(data, nil, prevEvent, t)
}

func TestParseEventWrongStart(t *testing.T) {

	//arrange
	testData := [][]byte{
		[]byte{},
		[]byte{vlq.Msb1Mask},
		[]byte{vlq.Msb1Mask, vlq.Msb1Mask},
		[]byte{vlq.Msb1Mask, vlq.Msb1Mask, vlq.MsbNo1Mask},
	}

	//tests
	for _, data := range testData {

		//act-assert
		EventAssert(data, nil, nil, t)
	}
}

func EventAssert(data []byte, realEvent smf.Event, prevEvent smf.Event, t *testing.T) {

	//act
	reader := bytes.NewBuffer(data)
	event, err := parseEvent(reader, prevEvent)

	//assert
	if realEvent == nil && err == nil {
		t.Error("Wait err but was nil")
	}

	if !reflect.DeepEqual(event, realEvent) && realEvent != nil {
		t.Errorf("Wait for event %v but was %v", realEvent, event)
	}
}
