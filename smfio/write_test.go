package smfio

import "testing"
import "github.com/algoGuy/EasyMIDI/smf"

func TestWriteNilWriter(t *testing.T) {

	//act
	err := Write(nil, &smf.MIDIFile{})

	//assert
	if err == nil {
		t.Error("wait error but was nil")
	}
}
