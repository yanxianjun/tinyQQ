package tlv

import "tinyQQ/lib/binary"

func T124(osType, osVersion, simInfo, apn []byte) []byte {
	return binary.NewWriterF(func(w *binary.Writer) {
		w.WriteUInt16(0x124)
		w.WriteTlv(binary.NewWriterF(func(w *binary.Writer) {
			w.WriteTlvLimitedSize(osType, 16)
			w.WriteTlvLimitedSize(osVersion, 16)
			w.WriteUInt16(2) // Network type wifi
			w.WriteTlvLimitedSize(simInfo, 16)
			w.WriteTlvLimitedSize([]byte{}, 16)
			w.WriteTlvLimitedSize(apn, 16)
		}))
	})
}
