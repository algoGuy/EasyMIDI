package smf

import (
	"testing"
)

func TestGetValue(t *testing.T) {

	//arrange
	event := &MIDIEvent{}

	track := &Track{}
	track.AddEvent(event)

	iterator := track.GetIterator()

	//act
	result := iterator.GetValue()

	//assert
	if result != nil {
		t.Error("wait error but was nil")
	}
}

func TestGetValueNoEvents(t *testing.T) {

	//arrange
	iterator := newTrackIterator(&Track{})

	//act
	result := iterator.GetValue()

	//assert
	if result != nil {
		t.Error("wait nil but was not nil")
	}
}

func TestMoveNextFirst(t *testing.T) {

	//arrange
	event := &MIDIEvent{}

	track := &Track{}
	track.AddEvent(event)

	iterator := track.GetIterator()

	//act
	result := iterator.MoveNext()

	//assert
	if !result {
		t.Error("Wait true but was false")
	}

	if iterator.GetValue() != event {
		t.Errorf("wait for %v but was %v", iterator.GetValue(), event)
	}
}

func TestMoveNext(t *testing.T) {

	//arrange
	event := &MIDIEvent{}
	event2 := &MIDIEvent{}

	track := &Track{}
	track.AddEvent(event)
	track.AddEvent(event2)

	iterator := track.GetIterator()

	//act
	iterator.MoveNext()
	result := iterator.MoveNext()

	//assert
	if !result {
		t.Error("Wait true but was false")
	}

	if iterator.GetValue() != event2 {
		t.Errorf("wait for %v but was %v", iterator.GetValue(), event2)
	}
}

func TestMoveNextNoElements(t *testing.T) {

	//arrange
	iterator := newTrackIterator(&Track{})

	//act
	result := iterator.MoveNext()

	//assert
	if result {
		t.Error("Wait false butwas true")
	}
}
