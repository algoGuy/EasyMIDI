package smfio

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

const midiExt string = ".mid"
const testDataFolder string = "./TestData/"

func TestReadWriteTest(t *testing.T) {

	//arrange
	allfiles, _ := ioutil.ReadDir(testDataFolder)
	var midFiles []os.FileInfo

	for _, file := range allfiles {
		if filepath.Ext(file.Name()) == midiExt {
			midFiles = append(midFiles, file)
		}
	}

	//tests
	for _, midiFile := range midFiles {

		//arrange
		file, err := os.Open(testDataFolder + midiFile.Name())
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		data, _ := ioutil.ReadAll(file)
		reader := bytes.NewReader(data)
		writer := &bytes.Buffer{}

		//act
		midi, readErr := Read(reader)
		writeErr := Write(writer, midi)

		//assert
		if readErr != nil {
			t.Errorf("Read midi error on %s file", midiFile.Name())
		}

		if writeErr != nil {
			t.Errorf("Write midi error on %s file", midiFile.Name())
		}

		if !bytes.Equal(data, writer.Bytes()) {
			t.Errorf("Error on %s file", midiFile.Name())
		}

		file.Close()
	}
}
