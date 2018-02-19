package smf

import (
	"reflect"
	"testing"
)

func TestTrackFromArray(t *testing.T) {

	//arrange
	eventsArray := []Event{}

	event, _ := NewMIDIEvent(0, NoteOffStatus, 0, 0x00, 0x01)
	eventsArray = append(eventsArray, event)

	event2, _ := NewSysexEvent(0, SysexDataStatus, []byte{0x01, 0x02})
	eventsArray = append(eventsArray, event2)

	//act
	track, err := TrackFromArray(eventsArray)

	//assert
	if err != nil {
		t.Error("wait error but was nil")
	}

	for i, val := 0, track.eventsList.Front(); val != nil; val = val.Next() {

		if !reflect.DeepEqual(val.Value, eventsArray[i]) {
			t.Errorf("wait for %v but was %v", val.Value, eventsArray[i])
		}

		i++
	}
}

func TestTrackFromArrayNilArray(t *testing.T) {

	//act
	_, err := TrackFromArray(nil)

	//assert
	if err == nil {
		t.Error("wait error but was nil")
	}
}

func TestTrackFromArrayNilEvent(t *testing.T) {

	//arrange
	eventsArray := []Event{nil}

	//act
	_, err := TrackFromArray(eventsArray)

	//assert
	if err == nil {
		t.Error("wait error but was nil")
	}
}

func TestGetIterator(t *testing.T) {

	//arrange
	track := &Track{}

	event, _ := NewMIDIEvent(0, NoteOnStatus, 0, 0x00, 0x01)
	track.AddEvent(event)

	//act
	iterator := track.GetIterator()

	//assert
	if iterator.trackRef != track {
		t.Errorf("wait for %v but was %v", iterator.trackRef, track)
	}
}

func TestAddEvent(t *testing.T) {

	//arrange
	event, _ := NewMetaEvent(0, MetaCuePoint, []byte{0x00})
	track := &Track{}

	//act
	track.AddEvent(event)

	//assert
	if reflect.DeepEqual(track.eventsList.Back(), track) {
		t.Errorf("wait for %v but was %v", track, track.eventsList.Back())
	}
}

func TestAddEventNil(t *testing.T) {

	//arrange
	track := &Track{}

	//act
	err := track.AddEvent(nil)

	//assert
	if err == nil {
		t.Error("wait error but was nil")
	}
}

func TestGetAllEvents(t *testing.T) {

	//arrange
	event, _ := NewMetaEvent(0x00, MetaCuePoint, []byte{0x00})

	track := &Track{}
	track.AddEvent(event)
	track.AddEvent(event)

	//act
	result := track.GetAllEvents()

	//assert
	realResult := []Event{event, event}

	for i := 0; i < len(result); i++ {

		if !reflect.DeepEqual(realResult[i], result[i]) {
			t.Errorf("wait for %v but was %v", realResult[i], result[i])
		}
	}
}

func TestRemoveAt(t *testing.T) {

	//arrange
	event, _ := NewMetaEvent(0, MetaCuePoint, []byte{0x00})
	event2, _ := NewMetaEvent(12, MetaEndOfTrack, []byte{0x00})

	track := &Track{}
	track.AddEvent(event)
	track.AddEvent(event2)

	//act
	track.RemoveAt(0)

	//assert
	if !reflect.DeepEqual(track.eventsList.Front().Value, event2) {
		t.Errorf("Wait %v but was %v", event2, track.eventsList.Front().Value)
	}

	if track.Len() != 1 {
		t.Errorf("Wait 1 but was %d", track.Len())
	}
}

func TestRemoveAtOutRange(t *testing.T) {

	//arrange
	event, _ := NewMetaEvent(0, MetaCuePoint, []byte{0x00})
	track := &Track{}

	track.AddEvent(event)

	//act
	err := track.RemoveAt(3)

	//assert
	if err == nil {
		t.Error("wait error but was nil")
	}
}
