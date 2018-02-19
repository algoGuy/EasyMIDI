package smf

import (
	"fmt"

	"github.com/algoGuy/EasyMIDI/vlq"
)

//MetaStatus for all Meta events
const MetaStatus uint8 = 0xFF

//Meta events types
const (
	MetaSequenceNumber    uint8 = 0x00
	MetaText              uint8 = 0x01
	MetaCopyrightNotice   uint8 = 0x02
	MetaSequenceTrackName uint8 = 0x03
	MetaInstrumentName    uint8 = 0x04
	MetaLyric             uint8 = 0x05
	MetaMarker            uint8 = 0x06
	MetaCuePoint          uint8 = 0x07
	MetaMIDIChannelPrefix uint8 = 0x20
	MetaMIDIPort          uint8 = 0x21
	MetaEndOfTrack        uint8 = 0x2F
	MetaSetTempo          uint8 = 0x51
	MetaSMTPEOffset       uint8 = 0x54
	MetaTimeSignature     uint8 = 0x58
	MetaKeySignature      uint8 = 0x59
	MetaSequencerSpecific uint8 = 0x7F
)

//MetaEvent struct
type MetaEvent struct {
	deltaTime uint32
	data      []byte
	metaType  uint8
}

//GetData return data from MetaEvent
func (mEvent *MetaEvent) GetData() []byte {
	return mEvent.data
}

//GetDTime return events delta time
func (mEvent *MetaEvent) GetDTime() uint32 {
	return mEvent.deltaTime
}

//GetStatus returns events status
func (*MetaEvent) GetStatus() uint8 {
	return MetaStatus
}

//GetMetaType get type for meta event
func (mEvent *MetaEvent) GetMetaType() uint8 {
	return mEvent.metaType
}

//NewMetaEvent Creates new meta event
func NewMetaEvent(deltaTime uint32, metaType uint8, data []byte) (*MetaEvent, error) {

	//check data
	if data == nil {
		return nil, &MidiError{"nil data reference"}
	}

	//check delta time
	if deltaTime > vlq.MaxNumber {
		return nil, &MidiError{"delta time is out of range"}
	}

	//check substatus
	if !CheckMetaType(metaType) {
		return nil, &MidiError{"wrong MetaEvent substatus"}
	}

	//check size
	if uint32(len(data)) > vlq.MaxNumber {
		return nil, &MidiError{"data size is out of vlq.MaxNumber range"}
	}

	//return new struct
	return &MetaEvent{
		deltaTime,
		data,
		metaType,
	}, nil
}

//CheckMetaStatus return true if status is MetaStatus
func CheckMetaStatus(status uint8) bool {
	return status == MetaStatus
}

//SetDtime sets deltaTime for event
func (mEvent *MetaEvent) SetDtime(deltaTime uint32) error {

	//check deltaTime
	if deltaTime > vlq.MaxNumber {
		return &MidiError{"deltaTime is out of range"}
	}

	mEvent.deltaTime = deltaTime

	return nil
}

//CheckMetaType return true if type is meta type
func CheckMetaType(metaType uint8) bool {

	//check 0x00 - 0x07 events
	if metaType <= MetaCuePoint {
		return true
	}

	//WHERE IS ENUMS? PLEASE!
	switch metaType {
	case MetaMIDIChannelPrefix, MetaEndOfTrack, MetaSetTempo, MetaSMTPEOffset,
		MetaTimeSignature, MetaKeySignature, MetaSequencerSpecific, MetaMIDIPort:
		return true
	default:
		return false
	}
}

//String return string representation for meta event
func (mEvent *MetaEvent) String() string {
	return fmt.Sprintf("META_EVENT: (deltaTime: %d metaType: %X data: % X)", mEvent.deltaTime, mEvent.metaType, mEvent.data)
}
