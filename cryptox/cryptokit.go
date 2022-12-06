package cryptox

import (
	"io"
	"os"
)

// MagicNumber is a defined 4-byte number to identify file type
// refer to https://en.wikipedia.org/wiki/List_of_file_signatures
// Header layout: magic number | key checksum
const MagicNumber = "\xFE\xF1\xFD\xEA"
const MagicNumberSize = len(MagicNumber)
const HeaderSize = MagicNumberSize + KeyHashSize
const encryptionBlockSize = 1 << 20

func IsEncrypted[S string | []byte](s S) bool {
	switch v := any(s).(type) {
	case string:
		sf, err := os.Open(v)
		if err != nil {
			return false
		}
		defer sf.Close()

		var header [HeaderSize]byte
		_, err = io.ReadFull(sf, header[:])
		if err != nil {
			return false
		}
		return string(header[:MagicNumberSize]) == MagicNumber
	case []byte:
		if len(v) < HeaderSize {
			return false
		}
		return string(v[:MagicNumberSize]) == MagicNumber
	default:
		return false
	}
}
