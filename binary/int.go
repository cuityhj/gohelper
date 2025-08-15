package binary

import "encoding/binary"

func Uint32ToBigEndianBytes(u32 uint32) []byte {
	fourBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(fourBytes, u32)
	return fourBytes
}

func Uint32ToLittleEndianBytes(u32 uint32) []byte {
	fourBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(fourBytes, u32)
	return fourBytes
}

func Uint64ToBigEndianBytes(u64 uint64) []byte {
	eightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(eightBytes, u64)
	return eightBytes
}

func Uint64ToLittleEndianBytes(u64 uint64) []byte {
	eightBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(eightBytes, u64)
	return eightBytes
}
