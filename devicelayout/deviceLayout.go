package devicelayout

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/go-playground/validator/v10"
)

type DeviceLayout struct {
	devDescriptors map[string]DeviceDescriptor
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
	slog.Info("Loading device layout configuration")
	// Open json file
	file, err := os.Open("./deviceLayout.json")
	if err != nil {
		slog.Error("Error opening file", "error", err)
		return err
	}
	defer file.Close()

	//Read json file
	data, err := io.ReadAll(file)
	if err != nil {
		slog.Error("Error reading file", "error", err)
		return err
	}

	// Parse json data
	devDesc := map[string]DeviceDescriptor{}
	err = json.Unmarshal([]byte(data), &devDesc)
	if err != nil {
		slog.Error("Error parsing file", "error", err)
		return err
	}

	if err = validate(devDesc); err != nil {
		return err
	}

	d.devDescriptors = devDesc

	return nil
}

func (d *DeviceLayout) FindLayout(vid, pid uint16) (*DeviceDescriptor, bool) {
	return d.FindLayoutByKey(fmt.Sprintf("%04x/%04x", vid, pid))
}

func (d *DeviceLayout) FindLayoutByKey(genKey string) (*DeviceDescriptor, bool) {
	dev, found := d.devDescriptors[genKey]
	return &dev, found
}

func validate(devDesc map[string]DeviceDescriptor) error {
	validate := validator.New()
	for key, desc := range devDesc {
		err := validate.Struct(desc)
		if err != nil {
			slog.Error("Error validating device layout configuration", "key", key)
			if _, ok := err.(*validator.InvalidValidationError); ok {
				slog.Error("", "error", err)
				return err
			}
			for _, err := range err.(validator.ValidationErrors) {
				slog.Error("Validation error", "field", err.Namespace(), "tag", err.Tag())
			}
			return err
		}
		if key != desc.Identifier.String() {
			err = fmt.Errorf("key does not match the VID and PID [%s != %s]", key, desc.Identifier.String())
			slog.Error("", "error", err)
			return err
		}
	}
	return nil
}
