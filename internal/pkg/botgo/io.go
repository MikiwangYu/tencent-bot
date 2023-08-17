package botgo

import (
	"bytes"
	"io"
)

func ReadAll(reader io.Reader, blockSize int, maxSize int) (data []byte, err error) {
	var blocks [][]byte
	for maxSize > 0 {
		var size int = min(blockSize, maxSize)
		var block []byte
		block = make([]byte, size, size)

		size, err = reader.Read(block)

		maxSize -= size
		block = block[:size]
		blocks = append(blocks, block)

		if err != nil {
			break
		}
	}

	data = bytes.Join(blocks, nil)
	return data, err
}
