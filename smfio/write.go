package smfio

import (
	"io"

	"bitbucket.org/NewStreeter/MIDIParser/smf"
)

//Write save MIDI to writer
func Write(writer io.Writer, midi *smf.MIDIFile) error {

	//check writer
	if writer == nil {
		return &WriteError{"nil writer reference"}
	}

	//write header data
	writeHeader(midi, writer)

	//write tracks + return
	return writeAllTracks(midi, writer)
}
