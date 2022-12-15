package util

const (
	SyntaxErr             = "ERR syntax error"
	InvalidIntErr         = "ERR value is not an integer or out of range"
	InvalidFloatErr       = "ERR value is not a valid float"
	InvalidLexErr         = "ERR min or max not valid string range item"
	WrongTypeErr          = "WRONGTYPE Operation against a key holding the wrong kind of value"
	WrongNumOfArgsErr     = "ERR wrong number of arguments for '%s' command"
	ZeroArgumentErr       = "ERR zero argument passed to the handler. This is an implementation bug"
	DeserializationErr    = "ERR unable to deserialize '%s' into a valid object"
	OptionNotSupportedErr = "ERR option '%s' is not currently supported"
	NegativeIntErr        = "ERR %s must be greater than 0"
	MustBePositiveErr     = "ERR %s must be positive"
)
