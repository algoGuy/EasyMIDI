package smf

import (
	"bytes"
	"testing"

	"github.com/algoGuy/EasyMIDI/vlq"
)

func TestNewSysexEvent(t *testing.T) {

	//arrange
	deltaTime := vlq.MaxNumber
	status := SysexStatus
	testData := []byte{0, 1, 2, 3, 4}

	//act
	result, err := NewSysexEvent(deltaTime, status, testData)

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

	if bytes.Compare(testData, result.GetData()) != 0 {
		t.Error("data not match")
	}
}

func TestNewSysexEventWrongDeltaTime(t *testing.T) {

	//act
	_, err := NewSysexEvent(vlq.MaxNumber+1, SysexStatus, []byte{})

	//assert
	if err == nil {
		t.Error("wrong delta time no error")
	}
}

func TestNewSysexNilData(t *testing.T) {

	//act
	_, err := NewSysexEvent(0, SysexStatus, nil)

	//assert
	if err == nil {
		t.Error("wrong data no error")
	}
}

func TestNewSysexWrongStatus(t *testing.T) {

	//act
	_, err := NewSysexEvent(0, 0x00, []byte{})

	//assert
	if err == nil {
		t.Error("wrong status no error")
	}
}

func TestNewSysexWrongDataSize(t *testing.T) {

	//act
	_, err := NewSysexEvent(0, SysexStatus, make([]byte, vlq.MaxNumber+1))

	//assert
	if err == nil {
		t.Error("wrong data size not nil")
	}
}

func TestCheckSysex(t *testing.T) {

	//arrange
	status1, status2 := SysexStatus, SysexDataStatus

	//act
	result1 := CheckSysexStatus(status1)
	result2 := CheckSysexStatus(status2)

	//assert
	if !result1 || !result2 {
		t.Error("Test case wait true but was false")
	}
}

func TestCheckSysexWrong(t *testing.T) {

	//arrange
	testStatus := SysexStatus - 1

	//act
	result := CheckSysexStatus(testStatus)

	//assert
	if result {
		t.Errorf("Test case %x wait false but was true", testStatus)
	}
}

func TestSetDtime(t *testing.T) {

	//arrange
	sEvent := &SysexEvent{}
	deltaTime := vlq.MaxNumber

	//act
	err := sEvent.SetDtime(deltaTime)

	//assert
	if sEvent.deltaTime != deltaTime {
		t.Error("Wrong delta time!")
	}

	if err != nil {
		t.Error("Must be nil return!")
	}
}

func TestSetDtimeOutOfRange(t *testing.T) {

	//arrange
	sEvent := &SysexEvent{}
	deltaTime := vlq.MaxNumber + 1

	//act
	err := sEvent.SetDtime(deltaTime)

	//assert
	if err == nil {
		t.Error("Delta time is out of range!", err)
	}
}
