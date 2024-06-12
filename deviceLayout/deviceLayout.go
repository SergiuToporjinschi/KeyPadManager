package devicelayout

import (
	"encoding/json"
	"fmt"
	"io"
	"main/logger"
	"main/utility"
	"os"
	"strconv"
	"sync"

	"github.com/go-playground/validator/v10"
)

type DeviceLayout struct {
	layouts map[string]DeviceLayoutConfig
}

var instance *DeviceLayout
var once sync.Once

func GetInstance() *DeviceLayout {
	once.Do(func() {
		instance = &DeviceLayout{}
	})
	return instance
}

func (d *DeviceLayout) LoadConfig() error {
	logger.Log.Info("Loading device layout configuration")
	// Open json file
	file, err := os.Open("./deviceLayout.json")
	if err != nil {
		logger.Log.Errorf("Error opening file: %v", err)
		return err
	}
	defer file.Close()

	//Read json file
	data, err := io.ReadAll(file)
	if err != nil {
		logger.Log.Errorf("Error reading file: %v", err)
		return err
	}

	// Parse json data
	layouts := map[string]DeviceLayoutConfig{}
	err = json.Unmarshal([]byte(data), &layouts)
	if err != nil {
		logger.Log.Errorf("Error parsing file: %v", err)
		return err
	}

	if err = validate(layouts); err != nil {
		return err
	}

	d.layouts = layouts

	return nil
}

func (d *DeviceLayout) FindLayout(vid, pid uint16) (*DeviceLayoutConfig, bool) {
	genKey := fmt.Sprintf("%04x/%04x", vid, pid)
	dev, found := d.layouts[genKey]
	return &dev, found
}
func (d *DeviceLayout) FindLayoutByKey(genKey string) (*DeviceLayoutConfig, bool) {
	dev, found := d.layouts[genKey]
	return &dev, found
}

func validate(layouts map[string]DeviceLayoutConfig) error {
	validate := validator.New()
	validate.RegisterValidation("byteNumber", validateByteNumber)
	for key, layout := range layouts {
		err := validate.Struct(layout)
		if err != nil {
			logger.Log.Errorf("Error validating device layout configuration for key %s: ", key)
			if _, ok := err.(*validator.InvalidValidationError); ok {
				logger.Log.Error(err)
				return err
			}
			for _, err := range err.(validator.ValidationErrors) {
				logger.Log.Errorf("Validation error: Field '%s' failed on the '%s' tag", err.Namespace(), err.Tag())
			}
			return err
		}
		if genKey := fmt.Sprintf("%04x/%04x", layout.Identifier.VID, layout.Identifier.PID); key != genKey {
			err = fmt.Errorf("key does not match the VID and PID [%s != %s]", key, genKey)
			logger.Log.Error(err)
			return err
		}
	}
	return nil
}

func validateByteNumber(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	if val == "" {
		return true
	}
	value, err := strconv.ParseUint(val, 2, 8)
	if err != nil {
		return false
	}

	// Check if the parsed value is within the valid byte range (0-255)
	return utility.NewIntSetWithValues(1, 2, 4, 8, 16, 32, 64, 128).Contains(int(value))
}
