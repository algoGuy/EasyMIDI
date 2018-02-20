# EasyMidi

[![codecov](https://codecov.io/gh/algoGuy/EasyMIDI/branch/master/graph/badge.svg)](https://codecov.io/gh/algoGuy/EasyMIDI) [![Build Status](https://semaphoreci.com/api/v1/algoguy/easymidi/branches/master/badge.svg)](https://semaphoreci.com/algoguy/easymidi) [![Go Report Card](https://goreportcard.com/badge/github.com/algoGuy/EasyMIDI)](https://goreportcard.com/report/github.com/algoGuy/EasyMIDI)

EasyMidi is a simple and reliable library for working with standard midi file (SMF).

## Installing

A step by step series of examples that tell you have to get a development env running

Get repository

```
go get github.com/algoGuy/easyMIDI
```

## How To Use
#### Example 1. Read and get data from midi file.
```go
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/algoGuy/smfio"
)

func main() {

	// Open test midi file
	file, _ := os.Open("./Test_-_test1.mid")
	defer file.Close()

	// Read and save midi to smf.MIDIFile struct
	midi, err := smfio.Read(bufio.NewReader(file))

	if err != nil {
		fmt.Println(err)
	}

	// Get zero track from
	track := midi.GetTrack(0)

	// Print number of midi tracks
	fmt.Println(midi.GetTracksNum())

	// Get all midi events via iterator
	iter := track.GetIterator()

	for iter.MoveNext() {
		fmt.Println(iter.GetValue())
	}
}
```
#### Example 2. Create and write one midi track into new midi file.
Created midi file from scratch. In current example output midi must contains one note.
```go
package main

import (
	"bufio"
	"log"
	"os"

	"github.com/algoGuy/smf"
	"github.com/algoGuy/smfio"
)

func main() {

	// Create division
	division, err := smf.NewDivision(960, smf.NOSMTPE)
	checkErr(err)

	// Create new midi struct
	midi, err := smf.NewSMF(smf.Format0, *division)
	checkErr(err)

	// Create track struct
	track := &smf.Track{}

	// Add track to new midi struct
	err = midi.AddTrack(track)
	checkErr(err)

	// Create some midi and meta events
	midiEventOne, err := smf.NewMIDIEvent(0, smf.NoteOnStatus, 0x00, 0x30, 0x50)
	checkErr(err)
	midiEventTwo, err := smf.NewMIDIEvent(10000, smf.NoteOnStatus, 0x00, 0x30, 0x00)
	checkErr(err)
	metaEventOne, err := smf.NewMetaEvent(0, smf.MetaEndOfTrack, []byte{})
	checkErr(err)

	// Add created events to track
	err = track.AddEvent(midiEventOne)
	checkErr(err)
	err = track.AddEvent(midiEventTwo)
	checkErr(err)
	err = track.AddEvent(metaEventOne)
	checkErr(err)

	// Save to new midi source file
	outputMidi, err := os.Create("outputMidi.mid")
	checkErr(err)
	defer outputMidi.Close()

	// Create buffering stream
	writer := bufio.NewWriter(outputMidi)
	smfio.Write(writer, midi)
	writer.Flush()
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
```

## Built With

* [Go](https://golang.org/) - The Go Programming Language

## Authors

* **algoGuy** - *main developer* - [algoGuy](https://github.com/algoGuy)
* **iqhater** - *main contributer* - [iqhater](https://github.com/iqhater)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
