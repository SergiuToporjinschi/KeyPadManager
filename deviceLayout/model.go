package devicelayout

import (
	"fmt"
	"strconv"
)

type DeviceLayoutConfig struct {
	Identifier Identifier  `json:"identifier" validate:"required"`
	Components []Component `json:"components" validate:"unique=Name,dive"`
}
type Identifier struct {
	Manufacturer string    `json:"manufacturer" validate:"required"`
	Name         string    `json:"name" validate:"required"`
	VID          HexUint16 `json:"vid" validate:"required"`
	PID          HexUint16 `json:"pid" validate:"required"`
}

type Component struct {
	Name        string `json:"name" validate:"required"`
	Icon        string `json:"icon" validate:"omitempty"`
	Type        string `json:"type" validate:"required,oneof=button rotaryEncoder"`
	ByteNumber  int    `json:"byteNumber" validate:"omitempty"`
	BitPosition string `json:"bitPosition" validate:"byteNumber,required_if=Type button"`
}
type HexUint16 uint16

func (h *HexUint16) UnmarshalJSON(b []byte) error {
	// Remove the quotes from the JSON string
	s := string(b)
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return fmt.Errorf("invalid hex string: %s", s)
	}
	s = s[1 : len(s)-1]

	// Convert the hex string to uint16
	v, err := strconv.ParseUint(s, 16, 16)
	if err != nil {
		return err
	}
	*h = HexUint16(v)
	return nil
}
