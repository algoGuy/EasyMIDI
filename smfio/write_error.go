package smfio

// WriteError error while MIDI writing
type WriteError struct {
	errorString string
}

// Error implements standart error interface
func (e *WriteError) Error() string {
	return "Write MIDI error: " + e.errorString
}
