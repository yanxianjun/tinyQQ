package tlv

import "tinyQQ/lib/binary"

func T52D(devInfo []byte) []byte {
	return binary.NewWriterF(func(w *binary.Writer) {
		w.WriteUInt16(0x52d)
		w.WriteTlv(binary.NewWriterF(func(w *binary.Writer) {
			w.Write(devInfo)
		}))
	})
}
