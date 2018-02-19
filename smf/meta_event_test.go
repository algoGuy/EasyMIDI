package smf

import "testing"
import "bitbucket.org/NewStreeter/MIDIParser/vlq"

import "bytes"

func TestCheckMetaType(t *testing.T) {

	//arrange
	testData := []byte{
		MetaSequenceNumber,
		MetaText,
		MetaCopyrightNotice,
		MetaSequenceTrackName,
		MetaInstrumentName,
		MetaLyric,
		MetaMarker,
		MetaCuePoint,
		MetaMIDIChannelPrefix,
		MetaMIDIPort,
		MetaEndOfTrack,
		MetaSetTempo,
		MetaSMTPEOffset,
		MetaTimeSignature,
		MetaKeySignature,
		MetaSequencerSpecific,
	}

	for _, metaType := range testData {

		//act
		result := CheckMetaType(metaType)

		//assert
		if !result {
			t.Errorf("Test case %x wait true but was false", metaType)
		}
	}
}

func TestCheckMetaTypeWrong(t *testing.T) {

	//arrange
	testType := MetaCuePoint + 1

	//act
	result := CheckMetaType(testType)

	//assert
	if result {
		t.Errorf("Test case %x wait false but was true", testType)
	}

}

func TestCheckMetaStatus(t *testing.T) {

	//act
	resultGood := CheckMetaStatus(MetaStatus)
	resultBad := CheckMetaStatus(MetaStatus - 1)

	//assert
	if !resultGood || resultBad {
		t.Errorf("wrong meta status")
	}

}

func TestNewMetaEvent(t *testing.T) {

	//arrange
	deltaTime := vlq.MaxNumber
	metaType := MetaCopyrightNotice
	data := []byte{1, 1}

	//act
	result, err := NewMetaEvent(deltaTime, metaType, data)

	//assert
	if err != nil {
		t.Errorf("%s", err)
	}

	if result.GetDTime() != deltaTime {
		t.Error("delta time not match")
	}

	if result.GetMetaType() != metaType {
		t.Error("metaType time not match")
	}

	if bytes.Compare(data, result.GetData()) != 0 {
		t.Error("data time not match")
	}

	if result.GetStatus() != MetaStatus {
		t.Error("wrong MetaStatus")
	}
}

func TestNewMetaEventNil(t *testing.T) {

	//act
	_, err := NewMetaEvent(0, MetaCopyrightNotice, nil)

	//assert
	if err == nil {
		t.Error("wait for error")
	}
}

func TestNewMetaEventOverDelay(t *testing.T) {

	//act
	_, err := NewMetaEvent(vlq.MaxNumber+1, MetaCopyrightNotice, []byte{})

	//assert
	if err == nil {
		t.Error("wait for error")
	}
}

func TestNewMetaEventWrongMetaType(t *testing.T) {

	//act
	_, err := NewMetaEvent(vlq.MaxNumber, MetaCuePoint+1, []byte{})

	//assert
	if err == nil {
		t.Error("wait for error")
	}
}

func TestNewMetaEventWrongDataSize(t *testing.T) {

	//act
	_, err := NewMetaEvent(vlq.MaxNumber, MetaCuePoint, make([]byte, vlq.MaxNumber+1))

	//assert
	if err == nil {
		t.Error("wait for error")
	}
}

func TestMetaSetDtime(t *testing.T) {

	//arrange
	mEvent := &MetaEvent{}
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

func TestMetaSetDtimeOutOfRange(t *testing.T) {

	//arrange
	mEvent := &MetaEvent{}
	deltaTime := vlq.MaxNumber + 1

	//act
	err := mEvent.SetDtime(deltaTime)

	//assert
	if err == nil {
		t.Error("Delta time is out of range!", err)
	}
}
