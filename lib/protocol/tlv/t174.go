package tlv

import "tinyQQ/lib/binary"

func T174(data []byte) []byte {
	return binary.NewWriterF(func(w *binary.Writer) {
		w.WriteUInt16(0x174)
		w.WriteTlv(data)
	})
}
