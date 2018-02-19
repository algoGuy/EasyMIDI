package smf

import (
	"bytes"
	"testing"

	"github.com/algoGuy/EasyMIDI/vlq"
)

func TestCheckMidiStatus(t *testing.T) {

	//arrange
	testData := []byte{
		NoteOffStatus,
		NoteOnStatus,
		PKPressureStatus,
		ControllerChangeStatus,
		ProgramChangeStatus,
		CKPressureStatus,
		PitchBendStatus,
	}

	for _, midiStatus := range testData {

		//act
		result := CheckMIDIStatus(midiStatus)

		//assert
		if !result {
			t.Errorf("Test case %x wait true but was false", midiStatus)
		}
	}
}

func TestCheckMidiWrongStatus(t *testing.T) {

	//arrange
	testStatus := NoteOffStatus - 1

	//act
	result := CheckMIDIStatus(testStatus)

	//assert
	if result {
		t.Errorf("Test case %x wait false but was true", testStatus)
	}

}

func TestNewMidiEvent(t *testing.T) {

	//arrange
	deltaTime := vlq.MaxNumber
	status := NoteOffStatus
	channel := MaxChannelNumber
	firstDataByte := MaxDataByteSize
	secondDataByte := MaxDataByteSize

	//act
	result, err := NewMIDIEvent(deltaTime, status, channel, firstDataByte, secondDataByte)

	//assert
	if err != nil {
		t.Error("unexpected error")
	}

	if result.GetDTime() != deltaTime {
		t.Error("delta time not match")
	}

	if result.GetStatus() != status {
		t.Error("status not match")
	}

	if bytes.Compare([]byte{firstDataByte, secondDataByte}, result.GetData()) != 0 {
		t.Error("data not match")
	}

	if result.GetChannel() != channel {
		t.Error("channel not match")
	}
}

func TestNewMidiEventOverDelay(t *testing.T) {

	//act
	_, err := NewMIDIEvent(vlq.MaxNumber+1, NoteOffStatus, 0x00, 0x00, 0x00)

	//assert
	if err == nil {
		t.Error("wait for error")
	}
}

func TestNewMidiEventWrongStatus(t *testing.T) {

	//act
	_, err := NewMIDIEvent(vlq.MaxNumber, 0x00, 0x00, 0x00, 0x00)

	//assert
	if err == nil {
		t.Error("wait for error")
	}
}

func TestNewMidiEventWrongChannel(t *testing.T) {

	//act
	_, err := NewMIDIEvent(vlq.MaxNumber, NoteOffStatus, MaxChannelNumber+1, 0x00, 0x00)

	//assert
	if err == nil {
		t.Error("wait for error")
	}
}

func TestNewMidiEventWrongData(t *testing.T) {

	//act
	_, err1 := NewMIDIEvent(vlq.MaxNumber, NoteOffStatus, MaxChannelNumber, MaxDataByteSize+1, 0x00)
	_, err2 := NewMIDIEvent(vlq.MaxNumber, NoteOffStatus, MaxChannelNumber, 0x00, MaxDataByteSize+1)

	//assert
	if err1 == nil || err2 == nil {
		t.Error("wait for error")
	}
}

func TestGetSingleData(t *testing.T) {

	//arrange
	firstDataByte := uint8(0x01)

	//act
	result, _ := NewMIDIEvent(0, ProgramChangeStatus, 0, firstDataByte, 0x00)

	//assert
	if bytes.Compare([]byte{firstDataByte}, result.GetData()) != 0 {
		t.Error("wrong data")
	}
}

func TestCheckSingleMIDIEvent(t *testing.T) {

	//act
	result1 := CheckSingleMIDIEvent(ProgramChangeStatus)
	result2 := CheckSingleMIDIEvent(CKPressureStatus)

	//assert
	if !result1 || !result2 {
		t.Error("expect true but was false")
	}
}

func TestCheckDoubleMIDIEvent(t *testing.T) {

	//act
	result := CheckSingleMIDIEvent(NoteOffStatus)

	//assert
	if result {
		t.Error("expect false but was true")
	}
}

func TestMidiSetDtime(t *testing.T) {

	//arrange
	mEvent := &MIDIEvent{}
	deltaTime := vlq.MaxNumber

	//act
	err := mEvent.SetDtime(deltaTime)

	//assert
	if mEvent.deltaTime != deltaTime {
		t.Error("Wrong delta time!")
	}

	if err != nil {
		t.Error("Must be nil return!")
	}
}

func TestMidiSetDtimeOutOfRange(t *testing.T) {

	//arrange
	mEvent := &MIDIEvent{}
	deltaTime := vlq.MaxNumber + 1

	//act
	err := mEvent.SetDtime(deltaTime)

	//assert
	if err == nil {
		t.Error("Delta time is out of range!", err)
	}
}
