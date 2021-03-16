package json

import "bytes"

func encodeWithIndent(dst *bytes.Buffer, src []byte, prefix, indentStr string) error {
	length := int64(len(src))
	indentNum := 0
	indentBytes := []byte(indentStr)
	for cursor := int64(0); cursor < length; cursor++ {
		c := src[cursor]
		switch c {
		case ' ', '\t', '\n', '\r':
			continue
		case '"':
			if err := dst.WriteByte(c); err != nil {
				return err
			}
			for {
				cursor++
				if err := dst.WriteByte(src[cursor]); err != nil {
					return err
				}
				switch src[cursor] {
				case '\\':
					cursor++
					if err := dst.WriteByte(src[cursor]); err != nil {
						return err
					}
				case '"':
					goto LOOP_END
				case nul:
					return errUnexpectedEndOfJSON("string", length)
				}
			}
		case '{':
			if cursor+1 < length && src[cursor+1] == '}' {
				if _, err := dst.Write([]byte{'{', '}'}); err != nil {
					return err
				}
				cursor++
			} else {
				indentNum++
				b := []byte{c, '\n'}
				b = append(b, prefix...)
				b = append(b, bytes.Repeat(indentBytes, indentNum)...)
				if _, err := dst.Write(b); err != nil {
					return err
				}
			}
		case '}':
			indentNum--
			if indentNum < 0 {
				return errInvalidCharacter('}', "}", cursor)
			}
			b := []byte{'\n'}
			b = append(b, prefix...)
			b = append(b, bytes.Repeat(indentBytes, indentNum)...)
			b = append(b, c)
			if _, err := dst.Write(b); err != nil {
				return err
			}
		case '[':
			if cursor+1 < length && src[cursor+1] == ']' {
				if _, err := dst.Write([]byte{'[', ']'}); err != nil {
					return err
				}
				cursor++
			} else {
				indentNum++
				b := []byte{c, '\n'}
				b = append(b, prefix...)
				b = append(b, bytes.Repeat(indentBytes, indentNum)...)
				if _, err := dst.Write(b); err != nil {
					return err
				}
			}
		case ']':
			indentNum--
			if indentNum < 0 {
				return errInvalidCharacter(']', "]", cursor)
			}
			b := []byte{'\n'}
			b = append(b, prefix...)
			b = append(b, bytes.Repeat(indentBytes, indentNum)...)
			b = append(b, c)
			if _, err := dst.Write(b); err != nil {
				return err
			}
		case ':':
			if _, err := dst.Write([]byte{':', ' '}); err != nil {
				return err
			}
		case ',':
			b := []byte{',', '\n'}
			b = append(b, prefix...)
			b = append(b, bytes.Repeat(indentBytes, indentNum)...)
			if _, err := dst.Write(b); err != nil {
				return err
			}
		default:
			if err := dst.WriteByte(c); err != nil {
				return err
			}
		}
	LOOP_END:
	}
	return nil
}
