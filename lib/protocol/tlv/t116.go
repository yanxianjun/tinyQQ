package tlv

import "tinyQQ/lib/binary"

func T116(miscBitmap, subSigMap uint32) []byte {
	return binary.NewWriterF(func(w *binary.Writer) {
		w.WriteUInt16(0x116)
		w.WriteTlv(binary.NewWriterF(func(w *binary.Writer) {
			w.WriteByte(0x00)
			w.WriteUInt32(miscBitmap)
			w.WriteUInt32(subSigMap)
			w.WriteByte(0x01)
			w.WriteUInt32(1600000226) // app id list
		}))
	})
}
