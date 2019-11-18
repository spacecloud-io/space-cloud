package utils

type LogLevel int

const(
	DEBUG       = iota
	INFO        = iota
	WARNING     = iota
	ERROR       = iota
)

func (level LogLevel) isValid() bool{
	return !(level < DEBUG || level > ERROR)
}

func (level LogLevel) String() string{
	names := [...]string{
		"Debug",
		"Info",
		"Warning",
		"Error" }

	if level < DEBUG || level > ERROR {
		return "Unknown"
	}

	return names[level]
}