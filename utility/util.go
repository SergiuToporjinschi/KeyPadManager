package utility

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/big"
	"regexp"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

func GetNoOfBitsAndBytes(value int64) (int, int) {
	// Example big integer
	bigNum := new(big.Int)
	bigNum.SetString(fmt.Sprintf("%d", value), 10)

	// Calculate the number of bits needed
	bitLen := bigNum.BitLen()

	// Calculate the number of bytes needed
	byteLen := (bitLen + 7) / 8
	return bitLen, byteLen
}

func FormatAsBinary(value, noOfBytes int) string {
	format := fmt.Sprintf("%%0%db", noOfBytes*8)
	result := fmt.Sprintf(format, value)

	re := regexp.MustCompile(".{1,8}")
	result = re.ReplaceAllStringFunc(result, func(s string) string {
		return s + " "
	})

	return fmt.Sprintf("[%s]", strings.TrimSpace(result))
}

func AbsInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func IntPointer(val int) *int {
	return &val
}

func NewIntBinding(val int) binding.ExternalInt {
	var value int = val
	return binding.BindInt(&value)
}

// RotateImageResource Rotate an image resource by a given angle
func RotateImageResource(res *fyne.StaticResource, angle float64) (*fyne.StaticResource, error) {
	srcImg, _, err := image.Decode(bytes.NewReader(res.StaticContent))
	if err != nil {
		return nil, err // Return error if image decoding fails
	}

	angle = angle * math.Pi / 180 // Convert angle to radians

	// Assuming the image is square for simplicity
	bounds := srcImg.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	dstImg := image.NewRGBA(image.Rect(0, 0, width, height))

	cx, cy := float64(width)/2, float64(height)/2

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dx, dy := float64(x)-cx, float64(y)-cy
			ox := cx + (dx*math.Cos(-angle) - dy*math.Sin(-angle))
			oy := cy + (dx*math.Sin(-angle) + dy*math.Cos(-angle))

			if ox >= 0 && ox < float64(width) && oy >= 0 && oy < float64(height) {
				dstImg.Set(x, y, srcImg.At(int(ox), int(oy)))
			} else {
				dstImg.Set(x, y, color.Transparent) // Fill with transparent color
			}
		}
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, dstImg) // Encode to JPEG, adjust as needed
	if err != nil {
		return nil, err // Return error if image encoding fails
	}

	// Create a new fyne.StaticResource for the rotated image
	rotatedRes := fyne.NewStaticResource(fmt.Sprintf("%s_rotated", res.Name()), buf.Bytes())
	return rotatedRes, nil
}
