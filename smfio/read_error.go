package smfio

// ReadError check all errors
type ReadError struct {
	errorString string
}

// Error implements standart error interface
func (e *ReadError) Error() string {
	return "Read MIDI error: " + e.errorString
}
