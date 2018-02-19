package smf

import "fmt"
import "bitbucket.org/NewStreeter/MIDIParser/vlq"

//MIDI Event status
const (
	NoteOffStatus          uint8 = 0x80
	NoteOnStatus           uint8 = 0x90
	PKPressureStatus       uint8 = 0xA0
	ControllerChangeStatus uint8 = 0xB0
	ProgramChangeStatus    uint8 = 0xC0
	CKPressureStatus       uint8 = 0xD0
	PitchBendStatus        uint8 = 0xE0
)

//MIDI Event sizes
const (
	MIDIDoubleEventSize uint32 = 3
	MIDISingleEventSize uint32 = 2
)

//MIDIEvent struct
type MIDIEvent struct {
	status         uint8
	deltaTime      uint32
	firstDataByte  uint8
	secondDataByte uint8
	channel        uint8
}

//MaxChannelNumber max number of channel
const MaxChannelNumber uint8 = 0x0F

//GetData gets midi event data
func (mEvent *MIDIEvent) GetData() []byte {

	//get result
	result := []byte{mEvent.firstDataByte, mEvent.secondDataByte}

	//check size of event
	if CheckSingleMIDIEvent(mEvent.status) {
		return result[:1]
	}

	//return result
	return result
}

//GetDTime get events Delta time
func (mEvent *MIDIEvent) GetDTime() uint32 {
	return mEvent.deltaTime
}

//GetStatus get events status
func (mEvent *MIDIEvent) GetStatus() uint8 {
	return mEvent.status
}

//GetChannel gets channel for MIDI event
func (mEvent *MIDIEvent) GetChannel() uint8 {
	return mEvent.channel
}

//SetDtime sets deltaTime for event
func (mEvent *MIDIEvent) SetDtime(deltaTime uint32) error {

	//check deltaTime
	if deltaTime > vlq.MaxNumber {
		return &MidiError{"deltaTime is out of range"}
	}

	mEvent.deltaTime = deltaTime

	return nil
}

//NewMIDIEvent create new midi event
func NewMIDIEvent(deltaTime uint32, status uint8, channel uint8, firstDataByte uint8, secondDataByte uint8) (*MIDIEvent, error) {

	//check deltaTime
	if deltaTime > vlq.MaxNumber {
		return nil, &MidiError{"deltaTime is out of range"}
	}

	//check status
	if !CheckMIDIStatus(status) {
		return nil, &MidiError{"Wrong MIDIEvent status!"}
	}

	//check channel number
	if channel > MaxChannelNumber {
		return nil, &MidiError{"Wrong channel number!"}
	}

	if firstDataByte > MaxDataByteSize || secondDataByte > MaxDataByteSize {
		return nil, &MidiError{"data bytes has msb = 1"}
	}

	//create struct
	return &MIDIEvent{
		status,
		deltaTime,
		firstDataByte,
		secondDataByte,
		channel,
	}, nil
}

//CheckMIDIStatus checks is status MIDIEvent Status
func CheckMIDIStatus(status uint8) bool {
	return status >= NoteOffStatus && status <= PitchBendStatus && status&0x0F == 0
}

//CheckSingleMIDIEvent return true for one byte size events
func CheckSingleMIDIEvent(status uint8) bool {
	return status == ProgramChangeStatus || status == CKPressureStatus
}

//String create string representation for MIDI event
func (mEvent *MIDIEvent) String() string {
	return fmt.Sprintf("MIDI_EVENT: (status: %X deltaTime: %d firstDataByte: %X secondDataByte: %X channel: %d)", mEvent.status, mEvent.deltaTime, mEvent.firstDataByte, mEvent.secondDataByte, mEvent.channel)
}
