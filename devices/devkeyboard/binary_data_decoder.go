package devkeyboard

import (
	"log/slog"
	"main/types"
)

func DecodeBinaryValue(data []byte, comp []DevKeyboardComponent) (types.IntSet, int) {
	if len(data) == 0 {
		return nil, 0
	}

	keyList := types.NewIntSet()
	encoderValue := 0

	for _, hrd := range comp {
		if hrd.Type == "button" {
			val := decodeButton(data[1:], hrd.BytesDescriptor)
			if val&hrd.Value != 0 {
				keyList.Add(hrd.Value)
			}
		} else if hrd.Type == "encoder" {
			encoderValue = decodeEncoderValue(data[1:], hrd.BytesDescriptor)
		}
	}
	return keyList, encoderValue
}

func decodeButton(data []byte, byteDesc BytesDescriptor) int {
	byteVal := make([]byte, byteDesc.Size)
	copy(byteVal, data[byteDesc.Start:byteDesc.Start+byteDesc.Size])
	var value int
	if byteDesc.Size > 1 { //more than one byte
		if byteDesc.Endianess == BYTE_DESCRIPTOR_ENDI_BIG {
			for i := 0; i < len(byteVal); i++ {
				value |= int(byteVal[i]) << (8 * (len(byteVal) - 1 - i))
			}
		} else if byteDesc.Endianess == BYTE_DESCRIPTOR_ENDI_LITTLE {
			for i := 0; i < len(byteVal); i++ {
				value |= int(byteVal[i]) << (8 * i)
			}
			byteTemp := byteVal[0]
			byteVal[0] = byteVal[1]
			byteVal[1] = byteTemp
		} else {
			slog.Warn("Endianess not specified")
		}
	} else { //one byte
		if byteDesc.Signed {
			value = int(int8(byteVal[0]))
		} else {
			value = int(uint8(byteVal[0]))
		}
	}
	return value
}

func decodeEncoderValue(data []byte, byteDesc BytesDescriptor) int {
	byteVal := make([]byte, byteDesc.Size)
	copy(byteVal, data[byteDesc.Start:byteDesc.Start+byteDesc.Size])

	if byteDesc.Signed {
		return int(int8(byteVal[0]))
	} else {
		return int(uint8(byteVal[0]))
	}
}
