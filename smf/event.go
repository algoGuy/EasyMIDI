package smf

//Event represent any type of midi events
type Event interface {
	SetDtime(uint32) error
	GetDTime() uint32
	GetStatus() uint8
	GetData() []byte
	String() string
}
