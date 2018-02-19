package smf

import (
	"math"
	"reflect"
	"testing"
)

func TestGetTrack(t *testing.T) {

	//arrange
	tracksNum := 7
	tracksSlice, midi := addTracksToSlice(tracksNum)

	//act-assert
	for i := 0; i < tracksNum; i++ {

		if tracksSlice[i] != midi.GetTrack(uint16(i)) {
			t.Error("Tracks is not equal!")
		}
	}
}

func addTracksToSlice(tracksNum int) ([]*Track, *MIDIFile) {

	midi := &MIDIFile{format: 1}
	tracksSlice := make([]*Track, tracksNum)

	for i := 0; i < tracksNum; i++ {
		track := &Track{}
		tracksSlice[i] = track
		midi.AddTrack(track)
	}
	return tracksSlice, midi
}

func TestGetTrackWrongNumber(t *testing.T) {

	//arrange
	tracksNum := 5
	_, midi := addTracksToSlice(tracksNum)

	//act
	result := midi.GetTrack(uint16(tracksNum + 1))

	//assert
	if result != nil {
		t.Error("Wrong number!")
	}
}

func TestRemoveTrack(t *testing.T) {

	//arrange
	tracksNum := 10
	trackToRemove := 5
	tracks, midi := addTracksToSlice(tracksNum)
	tracks = append(tracks[:trackToRemove], tracks[trackToRemove+1:]...)

	//act
	err := midi.RemoveTrack(uint16(trackToRemove))

	//assert
	for i := 0; i < len(tracks); i++ {

		if tracks[i] != midi.GetTrack(uint16(i)) {
			t.Error("Tracks is not equal!")
		}
	}

	if err != nil {
		t.Error("Expected nil but was error!")
	}
}

func TestRemoveTrackWrongNumber(t *testing.T) {

	//arrange
	tracksNum := 3
	trackToRemove := 5
	_, midi := addTracksToSlice(tracksNum)

	//act
	err := midi.RemoveTrack(uint16(trackToRemove))

	//assert
	if err == nil {
		t.Error("Wrong number to remove. Must be not nil error!")
	}
}

func GetTrackNum(t *testing.T) {

	//arrange
	tracksNum := 5
	_, midi := addTracksToSlice(tracksNum)

	//act
	result := midi.GetTracksNum()

	//assert
	if result != uint16(tracksNum) {
		t.Error("The numbers are not equal!")
	}
}

func TestAddTrackNil(t *testing.T) {

	//arrange
	midi := new(MIDIFile)

	//act
	err := midi.AddTrack(nil)

	//assert
	if err == nil {
		t.Error("Must be non nil return!")
	}
}

func TestAddTrackWrongFormat(t *testing.T) {

	//arrange
	midi := &MIDIFile{format: Format0}

	//act
	midi.AddTrack(&Track{})
	err := midi.AddTrack(&Track{})

	//assert
	if err == nil {
		t.Error("Must be non nil return!")
	}
}

func TestAddTrackMaxLimit(t *testing.T) {

	//arrange
	tracksNum := math.MaxUint16
	_, midi := addTracksToSlice(tracksNum)

	//act
	err := midi.AddTrack(&Track{})

	//assert
	if err == nil {
		t.Error("Must be non nil return!")
	}
}

/* func TestGetFormat(t *testing.T) {

	//arrange
	midi := &MIDIFile{format: Format1}

	//act
	result := midi.format

	//assert
	if result != midi.GetFormat() {
		t.Error("Wrong midi format!")
	}
} */

func TestNewMIDI(t *testing.T) {

	//arrange
	format := Format2
	division, _ := NewDivision(960, NOSMTPE)

	//act
	midi, err := NewSMF(format, *division)

	//assert
	if err != nil {
		t.Error("Must be nil return!")
	}

	if midi.GetFormat() != format {
		t.Error("Wrong format!")
	}

	if !reflect.DeepEqual(midi.GetDivision(), *division) {
		t.Error("Wrong division!")
	}
}

func TestNewMIDIWrongFormat(t *testing.T) {

	//arrange
	format := Format2 + 1
	division := Division{}

	//act
	result, _ := NewSMF(format, division)

	//assert
	if result != nil {
		t.Error("Wrong format! Must be nil return!")
	}
}
