package vlq

// Error check all errors
type Error struct {
	errorString string
}

// Error implements standart error interface
func (e *Error) Error() string {
	return "Vlq error: " + e.errorString
}
