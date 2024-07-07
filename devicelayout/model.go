package devicelayout

import (
	"fmt"
	"strconv"
)

type (
	DeviceDescriptor struct {
		Identifier         Identifier `json:"identifier" validate:"required"`
		Report             Report     `json:"report" validate:"required"`
		HardwareDescriptor any        `json:"hardwareDescriptor"`
	}
	Identifier struct {
		Manufacturer string    `json:"manufacturer" validate:"required"`
		Product      string    `json:"product" validate:"required"`
		SerialNumber string    `json:"serialNumber" validate:"omitempty"`
		VID          HexUint16 `json:"vid" validate:"required"`
		PID          HexUint16 `json:"pid" validate:"required"`
	}

	Report struct {
		Size uint8 `json:"size" validate:"required"`
	}

	HexUint16 uint16
)

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
func (h HexUint16) String() string {
	return strconv.FormatUint(uint64(h), 16)
}

func (di *Identifier) String() string {
	return fmt.Sprintf("%s/%s", di.VID, di.PID)
}

func (di *Identifier) StringDetailed() string {
	return fmt.Sprintf("%s by %s [%s]", di.Product, di.Manufacturer, di.String())
}
