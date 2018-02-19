package smf

//SMTPE time division midi valid codes
const (
	NOSMTPE int8 = 0
	SMTPE24 int8 = -24
	SMTPE25 int8 = -25
	SMTPE29 int8 = -29
	SMTPE30 int8 = -30
)

const (

	//TicksMaxValue midi division ticks max value with NOSMTPE
	TicksMaxValue uint16 = 0x7FFF

	//SMTPETicksMaxValue midi division ticks max value with SMTPE
	SMTPETicksMaxValue uint16 = 0xFF
)

//Division represents time division struct
type Division struct {
	ticks       uint16
	smtpeFrames int8
}

//GetTicks return ticks for division
func (d *Division) GetTicks() uint16 {
	return d.ticks
}

//GetSMTPE return SMTPE for division
func (d *Division) GetSMTPE() int8 {
	return d.smtpeFrames
}

//IsSMTPE checks is current division SMTPE
func (d *Division) IsSMTPE() bool {
	return d.smtpeFrames != NOSMTPE
}

//NewDivision creates new time division
func NewDivision(ticks uint16, SMTPEFrames int8) (*Division, error) {

	//check for SMTPE
	if !CheckSMTPE(SMTPEFrames) {
		return nil, &MidiError{"No supported SMTPE"}
	}

	//check for tics size
	if SMTPEFrames != NOSMTPE && ticks > SMTPETicksMaxValue {
		return nil, &MidiError{"Wrong ticks for SMTPE"}
	}

	//check msb for ticks
	if ticks > TicksMaxValue {
		return nil, &MidiError{"ticks has msb = 1"}
	}

	return &Division{ticks: ticks, smtpeFrames: SMTPEFrames}, nil
}

//CheckSMTPE return true if SMTPEFrames is valid values
func CheckSMTPE(SMTPEFrames int8) bool {

	switch SMTPEFrames {
	case NOSMTPE, SMTPE24, SMTPE25, SMTPE29, SMTPE30:
		return true
	default:
		return false
	}
}
