package smf

import (
	"fmt"

	"bitbucket.org/NewStreeter/MIDIParser/vlq"
)

// SysexStatus and SysexDataStatus is two flavors of sysex_events
const (
	SysexStatus     uint8 = 0xF0
	SysexDataStatus uint8 = 0xF7
)

//SysexEvent for midi
type SysexEvent struct {
	status    uint8
	deltaTime uint32
	data      []byte
}

//GetDTime get delta time for event
func (sEvent *SysexEvent) GetDTime() uint32 {
	return sEvent.deltaTime
}

//GetStatus get status for event
func (sEvent *SysexEvent) GetStatus() uint8 {
	return sEvent.status
}

//GetData get data for event
func (sEvent *SysexEvent) GetData() []byte {
	return sEvent.data
}

//SetDtime sets deltaTime for event
func (sEvent *SysexEvent) SetDtime(deltaTime uint32) error {

	//check deltaTime
	if deltaTime > vlq.MaxNumber {
		return &MidiError{"deltaTime is out of range"}
	}

	sEvent.deltaTime = deltaTime

	return nil
}

//NewSysexEvent creataes new SysexEvent
func NewSysexEvent(deltaTime uint32, status uint8, data []byte) (*SysexEvent, error) {

	// check delta time
	if deltaTime > vlq.MaxNumber {
		return nil, &MidiError{"deltaTime out of range"}
	}

	//check data
	if data == nil {
		return nil, &MidiError{"nil data reference"}
	}

	//check sysex
	if !CheckSysexStatus(status) {
		return nil, &MidiError{"Wrong SysexEvent status!"}
	}

	//check size
	if uint32(len(data)) > vlq.MaxNumber {
		return nil, &MidiError{"data size is out of vlq.MaxNumber range"}
	}

	return &SysexEvent{deltaTime: deltaTime, status: status, data: data}, nil
}

//CheckSysexStatus check is status SysexStatus
func CheckSysexStatus(status uint8) bool {
	return SysexStatus == status || SysexDataStatus == status
}

//String returns string representation of event
func (sEvent *SysexEvent) String() string {
	return fmt.Sprintf("SYSEX_EVENT: (status: %X deltaTime: %d data: % X)", sEvent.status, sEvent.deltaTime, sEvent.data)
}
