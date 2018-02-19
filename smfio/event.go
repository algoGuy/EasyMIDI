package smfio

import (
	"encoding/binary"
	"io"

	"bitbucket.org/NewStreeter/MIDIParser/smf"
	"bitbucket.org/NewStreeter/MIDIParser/vlq"
)

//parseEvent parse next events from reader
func parseEvent(reader io.Reader, prevEvent smf.Event) (smf.Event, error) {

	//get event dtime
	dTime, err := vlq.GetFromReader(reader)
	if err != nil {
		return nil, &ReadError{"Track data corrupted"}
	}

	//get status
	var status byte
	err = binary.Read(reader, binary.BigEndian, &status)
	if err != nil {
		return nil, &ReadError{"Track data corrupted"}
	}

	//parse running status
	if status < vlq.Msb1Mask && prevEvent != nil && smf.CheckMIDIStatus(prevEvent.GetStatus()) {
		return parseRunningStatus(status, dTime, reader, prevEvent)
	}

	//parse midi event
	if smf.CheckMIDIStatus(status & 0xF0) {
		return parseMIDIEvent(status, dTime, reader)
	}

	//parse sysex event
	if smf.CheckSysexStatus(status) {
		return parseSysexEvent(status, dTime, reader)
	}

	//parse meta event
	if smf.CheckMetaStatus(status) {
		return parseMetaEvent(status, dTime, reader)
	}

	//not supported event - crash
	return nil, &ReadError{"Track data corrupted"}
}

//parseRunningStatus used for parse running status
func parseRunningStatus(status byte, dTime uint32, reader io.Reader, prevEvent smf.Event) (smf.Event, error) {

	//cast event
	midiEvent := prevEvent.(*smf.MIDIEvent)

	//check double byte midi event status
	var secondByte byte
	if !smf.CheckSingleMIDIEvent(prevEvent.GetStatus()) {

		//read second byte
		if err := binary.Read(reader, binary.BigEndian, &secondByte); err != nil {
			return nil, err
		}
	}

	//return event
	return smf.NewMIDIEvent(dTime, midiEvent.GetStatus(), midiEvent.GetChannel(), status, secondByte)
}

//parseMIDIEvent used for parse MIDI events
func parseMIDIEvent(status byte, dTime uint32, reader io.Reader) (smf.Event, error) {

	//calculate status & channel
	channel := status & 0x0F
	status &= 0xF0

	//read first byte
	var firstByte byte
	if err := binary.Read(reader, binary.BigEndian, &firstByte); err != nil {
		return nil, err
	}

	//read second byte
	var secondByte byte
	if !smf.CheckSingleMIDIEvent(status) {

		if err := binary.Read(reader, binary.BigEndian, &secondByte); err != nil {
			return nil, err
		}
	}

	return smf.NewMIDIEvent(dTime, status, channel, firstByte, secondByte)
}

//parseSysexEvent used for parse Sysex events only
func parseSysexEvent(status byte, deltaTime uint32, reader io.Reader) (smf.Event, error) {

	//get length
	length, err := vlq.GetFromReader(reader)
	if err != nil {
		return nil, err
	}

	//read data
	data := make([]byte, length)
	if _, err := io.ReadAtLeast(reader, data, int(length)); err != nil {
		return nil, &ReadError{"Sysex event data out of range"}
	}

	return smf.NewSysexEvent(deltaTime, status, data)
}

//parseMetaEvent used for parse Meta events only
func parseMetaEvent(status byte, deltaTime uint32, reader io.Reader) (smf.Event, error) {

	//get meta type
	var metaType byte
	if err := binary.Read(reader, binary.BigEndian, &metaType); err != nil {
		return nil, err
	}

	//check type
	if !smf.CheckMetaType(metaType) {
		return nil, &ReadError{"Wrong meta type for meta event"}
	}

	//get length
	length, err := vlq.GetFromReader(reader)
	if err != nil {
		return nil, err
	}

	//read data
	data := make([]byte, length)
	if _, err := io.ReadAtLeast(reader, data, int(length)); err != nil {
		return nil, &ReadError{"Meta event data out of range"}
	}

	return smf.NewMetaEvent(deltaTime, metaType, data)
}

//writeEvent writes event to writer
func writeEvent(event smf.Event, prevEvent smf.Event, writer io.Writer) error {

	//add vlq length
	writer.Write(vlq.GetBytes(event.GetDTime()))

	//get status and event data
	status := []byte{event.GetStatus()}
	eventData := event.GetData()

	//write events
	switch castEvent := event.(type) {
	case *smf.MIDIEvent:
		{
			//try cast prev event
			prevMidiEvent, ok := prevEvent.(*smf.MIDIEvent)

			//check running status
			if !ok || !(prevEvent.GetStatus() == status[0] && castEvent.GetChannel() == prevMidiEvent.GetChannel()) {
				writer.Write([]byte{status[0] | castEvent.GetChannel()})
			}
		}
	case *smf.MetaEvent:
		{
			writer.Write(status)
			writer.Write([]byte{castEvent.GetMetaType()})
			writer.Write(vlq.GetBytes(uint32(len(eventData))))
		}
	case *smf.SysexEvent:
		{
			writer.Write(status)
			writer.Write(vlq.GetBytes(uint32(len(eventData))))
		}
	default:
		{
			return &WriteError{"event type: " + event.String() + "not supported"}
		}
	}

	//write event data
	writer.Write(eventData)

	return nil
}
