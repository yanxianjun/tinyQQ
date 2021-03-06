package tlv

import "tinyQQ/lib/binary"

func T18(appId uint32, uin uint32) []byte {
	return binary.NewWriterF(func(w *binary.Writer) {
		w.WriteUInt16(0x18)
		w.WriteTlv(binary.NewWriterF(func(w *binary.Writer) {
			w.WriteUInt16(1)
			w.WriteUInt32(1536)
			w.WriteUInt32(appId)
			w.WriteUInt32(0)
			w.WriteUInt32(uin)
			w.WriteUInt16(0)
			w.WriteUInt16(0)
		}))
	})
}
