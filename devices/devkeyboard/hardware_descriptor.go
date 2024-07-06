package devkeyboard

import (
	"encoding/json"
	"log/slog"
	"main/devicelayout"
)

type DevKeyboardComponent struct {
	Name            string          `json:"name" validate:"required"`
	Type            string          `json:"type" validate:"required,oneof=button encoder encoderButton"`
	BytesDescriptor BytesDescriptor `json:"bytesDescriptor" validate:"required"`
	Value           int             `json:"value" validate:"required_if=Type button"`
}

type BytesDescriptor struct {
	Start     int    `json:"start"`
	Size      int    `json:"size"`
	Endianess string `json:"endianess" validate:"omitempty,oneof=big little"`
	Signed    bool   `json:"signed" validate:"omitempty"`
}

const (
	BYTE_DESCRIPTOR_ENDI_BIG    = "big"
	BYTE_DESCRIPTOR_ENDI_LITTLE = "little"
)

func ConvertHardwareDescriptor(devDescriptor *devicelayout.DeviceDescriptor) *devicelayout.DeviceDescriptor {
	var comps []DevKeyboardComponent

	desc, err := json.Marshal(devDescriptor.HardwareDescriptor)
	if err != nil {
		slog.Error("Error parsing devDescriptor.HardwareDescriptor", "error", err)
		panic(err)
	}

	err = json.Unmarshal(desc, &comps)
	if err != nil {
		slog.Error("Error parsing []DevKeyboardComponent", "error", err)
		panic(err)
	}

	devDescriptor.HardwareDescriptor = comps
	return devDescriptor
}

func (inst *DevKeyboardComponent) GetButtonValue() int {
	return inst.Value
}
