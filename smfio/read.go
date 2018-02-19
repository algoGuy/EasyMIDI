package smfio

import (
	"io"

	"bitbucket.org/NewStreeter/MIDIParser/smf"
)

const minTrackEndSize = 4

//Read reads midi from io.Reader
func Read(reader io.Reader) (*smf.MIDIFile, error) {

	//check reader
	if reader == nil {
		return nil, &ReadError{"nil Reader reference"}
	}

	//parse header
	header, err := parseHeader(reader)
	if err != nil {
		return nil, err
	}

	//create midi
	midi, err := smf.NewSMF(header.format, *header.division)
	if err != nil {
		return nil, &ReadError{err.Error()}
	}

	//parse tracks
	return parseAllTracks(midi, reader, header.tracksNum)
}
