package tlv

import (
	"tinyQQ/lib/binary"
)

func T100(ssoVersion, protocol, mainSigMap uint32) []byte {
	return binary.NewWriterF(func(w *binary.Writer) {
		w.WriteUInt16(0x100)
		w.WriteTlv(binary.NewWriterF(func(w *binary.Writer) {
			w.WriteUInt16(1)
			w.WriteUInt32(ssoVersion)
			w.WriteUInt32(16)
			w.WriteUInt32(protocol)
			w.WriteUInt32(0)          // App client version
			w.WriteUInt32(mainSigMap) // 34869472
		}))
	})
}
