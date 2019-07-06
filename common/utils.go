package common

import (
	"encoding/binary"
)

func EncodePacket(idx uint8, content []byte) []byte {
	head := []byte{idx}
	size := make([]byte, 2)
	binary.LittleEndian.PutUint16(size, uint16(len(content)))
	return append(append(head[:], size[:]...), content[:]...)
}

func DecodePacket(data []byte) ([]byte, uint8, []byte) {
	head := uint8(data[0])
	size := binary.LittleEndian.Uint16(data[1:3])

	if len(data) >= int(3 + size) {
		content := data[3:3+size]
		return data[3+size:], head, content
	}
	return data, 0, nil
}

func ConvertInt32ToBytes(i int32) []byte {
    var buffer = make([]byte, 4)
    binary.LittleEndian.PutUint32(buffer, uint32(i))
    return buffer
}

func ConvertBytesToInt32(buffer []byte) int32 {
    return int32(binary.LittleEndian.Uint32(buffer))
}


