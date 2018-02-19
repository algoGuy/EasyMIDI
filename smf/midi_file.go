package smf

import (
	"math"
)

//MIDI file formats
const (
	Format0 uint16 = 0
	Format1 uint16 = 1
	Format2 uint16 = 2
)

//MaxDataByteSize max size for event data
const MaxDataByteSize uint8 = 0x7F

//MIDIFile file struct
type MIDIFile struct {
	format      uint16
	division    Division
	tracksSlice []*Track
}

//GetTrack by number return nil if number > track num
func (midi *MIDIFile) GetTrack(number uint16) *Track {

	if number < uint16(len(midi.tracksSlice)) {
		return midi.tracksSlice[number]
	}

	return nil
}

//RemoveTrack remove track from midi
func (midi *MIDIFile) RemoveTrack(number uint16) error {

	//check number
	if number >= uint16(len(midi.tracksSlice)) {
		return &MidiError{"track number is out of range"}
	}

	//remove tack
	midi.tracksSlice = append(midi.tracksSlice[:number], midi.tracksSlice[number+1:]...)

	return nil
}

//GetTracksNum current number of tracks
func (midi *MIDIFile) GetTracksNum() uint16 {
	return uint16(len(midi.tracksSlice))
}

//AddTrack adds track to midi
func (midi *MIDIFile) AddTrack(track *Track) error {

	//nil check
	if track == nil {
		return &MidiError{"try to add nil track"}
	}

	//check format
	if midi.format == Format0 && len(midi.tracksSlice) != 0 {
		return &MidiError{"try to add multiple track"}
	}

	//check max track
	if midi.GetTracksNum() == math.MaxUint16 {
		return &MidiError{"adding track over MaxUint16 limit"}
	}

	//append new value
	midi.tracksSlice = append(midi.tracksSlice, track)

	return nil
}

//GetFormat return format for midi file
func (midi *MIDIFile) GetFormat() uint16 {
	return midi.format
}

//GetDivision return division for midi file
func (midi *MIDIFile) GetDivision() Division {
	return midi.division
}

//NewSMF create new midi file
func NewSMF(format uint16, division Division) (*MIDIFile, error) {

	//check format
	if format > Format2 {
		return nil, &MidiError{"Wrong midi format!"}
	}

	return &MIDIFile{format: format, division: division}, nil
}
