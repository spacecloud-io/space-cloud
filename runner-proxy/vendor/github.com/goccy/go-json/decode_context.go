package json

import (
	"unsafe"
)

var (
	isWhiteSpace = [256]bool{}
)

func init() {
	isWhiteSpace[' '] = true
	isWhiteSpace['\n'] = true
	isWhiteSpace['\t'] = true
	isWhiteSpace['\r'] = true
}

func char(ptr unsafe.Pointer, offset int64) byte {
	return *(*byte)(unsafe.Pointer(uintptr(ptr) + uintptr(offset)))
}

func skipWhiteSpace(buf []byte, cursor int64) int64 {
LOOP:
	if isWhiteSpace[buf[cursor]] {
		cursor++
		goto LOOP
	}
	return cursor
}

func skipObject(buf []byte, cursor, depth int64) (int64, error) {
	braceCount := 1
	for {
		switch buf[cursor] {
		case '{':
			braceCount++
			depth++
			if depth > maxDecodeNestingDepth {
				return 0, errExceededMaxDepth(buf[cursor], cursor)
			}
		case '}':
			depth--
			braceCount--
			if braceCount == 0 {
				return cursor + 1, nil
			}
		case '[':
			depth++
			if depth > maxDecodeNestingDepth {
				return 0, errExceededMaxDepth(buf[cursor], cursor)
			}
		case ']':
			depth--
		case '"':
			for {
				cursor++
				switch buf[cursor] {
				case '"':
					if buf[cursor-1] == '\\' {
						continue
					}
					goto SWITCH_OUT
				case nul:
					return 0, errUnexpectedEndOfJSON("string of object", cursor)
				}
			}
		case nul:
			return 0, errUnexpectedEndOfJSON("object of object", cursor)
		}
	SWITCH_OUT:
		cursor++
	}
}

func skipArray(buf []byte, cursor, depth int64) (int64, error) {
	bracketCount := 1
	for {
		switch buf[cursor] {
		case '[':
			bracketCount++
			depth++
			if depth > maxDecodeNestingDepth {
				return 0, errExceededMaxDepth(buf[cursor], cursor)
			}
		case ']':
			bracketCount--
			depth--
			if bracketCount == 0 {
				return cursor + 1, nil
			}
		case '{':
			depth++
			if depth > maxDecodeNestingDepth {
				return 0, errExceededMaxDepth(buf[cursor], cursor)
			}
		case '}':
			depth--
		case '"':
			for {
				cursor++
				switch buf[cursor] {
				case '"':
					if buf[cursor-1] == '\\' {
						continue
					}
					goto SWITCH_OUT
				case nul:
					return 0, errUnexpectedEndOfJSON("string of object", cursor)
				}
			}
		case nul:
			return 0, errUnexpectedEndOfJSON("array of object", cursor)
		}
	SWITCH_OUT:
		cursor++
	}
}

func skipValue(buf []byte, cursor, depth int64) (int64, error) {
	for {
		switch buf[cursor] {
		case ' ', '\t', '\n', '\r':
			cursor++
			continue
		case '{':
			return skipObject(buf, cursor+1, depth+1)
		case '[':
			return skipArray(buf, cursor+1, depth+1)
		case '"':
			for {
				cursor++
				switch buf[cursor] {
				case '"':
					if buf[cursor-1] == '\\' {
						continue
					}
					return cursor + 1, nil
				case nul:
					return 0, errUnexpectedEndOfJSON("string of object", cursor)
				}
			}
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			for {
				cursor++
				if floatTable[buf[cursor]] {
					continue
				}
				break
			}
			return cursor, nil
		case 't':
			buflen := int64(len(buf))
			if cursor+3 >= buflen {
				return 0, errUnexpectedEndOfJSON("bool of object", cursor)
			}
			if buf[cursor+1] != 'r' {
				return 0, errUnexpectedEndOfJSON("bool of object", cursor)
			}
			if buf[cursor+2] != 'u' {
				return 0, errUnexpectedEndOfJSON("bool of object", cursor)
			}
			if buf[cursor+3] != 'e' {
				return 0, errUnexpectedEndOfJSON("bool of object", cursor)
			}
			cursor += 4
			return cursor, nil
		case 'f':
			buflen := int64(len(buf))
			if cursor+4 >= buflen {
				return 0, errUnexpectedEndOfJSON("bool of object", cursor)
			}
			if buf[cursor+1] != 'a' {
				return 0, errUnexpectedEndOfJSON("bool of object", cursor)
			}
			if buf[cursor+2] != 'l' {
				return 0, errUnexpectedEndOfJSON("bool of object", cursor)
			}
			if buf[cursor+3] != 's' {
				return 0, errUnexpectedEndOfJSON("bool of object", cursor)
			}
			if buf[cursor+4] != 'e' {
				return 0, errUnexpectedEndOfJSON("bool of object", cursor)
			}
			cursor += 5
			return cursor, nil
		case 'n':
			buflen := int64(len(buf))
			if cursor+3 >= buflen {
				return 0, errUnexpectedEndOfJSON("null", cursor)
			}
			if buf[cursor+1] != 'u' {
				return 0, errUnexpectedEndOfJSON("null", cursor)
			}
			if buf[cursor+2] != 'l' {
				return 0, errUnexpectedEndOfJSON("null", cursor)
			}
			if buf[cursor+3] != 'l' {
				return 0, errUnexpectedEndOfJSON("null", cursor)
			}
			cursor += 4
			return cursor, nil
		default:
			return cursor, errUnexpectedEndOfJSON("null", cursor)
		}
	}
}
