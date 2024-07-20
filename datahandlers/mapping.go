package datahandlers

import (
	"encoding/json"
	"log/slog"
	"os"
	"sync"
)

const (
	mappingFileName = "mapping.json"
)

type MappingHandler[T any] struct {
	//deviceIdentifier -> appIdentifier -> key -> action
	List             map[string]map[string]map[int]Action[T] `json:"list"`
	mappingListMutex sync.Mutex                              `json:"-"`
}

type Action[T any] struct {
	Type  string
	Value T
}

var onceMappingHandler sync.Once
var instanceMappingHandler *MappingHandler[any]

func GetMappingHandlerInstance() *MappingHandler[any] {
	onceMappingHandler.Do(func() {
		instanceMappingHandler = &MappingHandler[any]{
			List: make(map[string]map[string]map[int]Action[any]),
		}
	})

	return instanceMappingHandler
}

func (mh *MappingHandler[T]) AddDeviceKeyMacro(devIdentifier string, appIdentifier string, devKey int, val T) {
	mh.mappingListMutex.Lock()
	defer mh.mappingListMutex.Unlock()

	mh.add("Macro", devIdentifier, appIdentifier, devKey, val)
}

func (mh *MappingHandler[T]) AddDeviceKeyKey(devIdentifier string, appIdentifier string, devKey int, val T) {
	mh.mappingListMutex.Lock()
	defer mh.mappingListMutex.Unlock()

	mh.add("Key", devIdentifier, appIdentifier, devKey, val)
}

func (mh *MappingHandler[T]) SaveMapping() {
	mh.mappingListMutex.Lock()
	defer mh.mappingListMutex.Unlock()

	jsonMappings, err := json.MarshalIndent(mh.List, "", "  ")
	if err != nil {
		slog.Error("Error marshalling app list: ", "error", err)
		return
	}

	file, err := os.OpenFile(mappingFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		slog.Error("Error opening file: ", "error", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(string(jsonMappings))
	if err != nil {
		slog.Error("Error writing to file: ", "error", err)
		return
	}
}

func (mh *MappingHandler[T]) LoadMapping() int {
	mh.ClearAll()

	mh.mappingListMutex.Lock()
	defer mh.mappingListMutex.Unlock()

	content, err := os.ReadFile(mappingFileName)
	if os.IsNotExist(err) {
		slog.Warn("File not found: ", "error", err)
		return 0
	} else if err != nil {
		slog.Error("Error opening file: ", "error", err)
	}

	var mapping map[string]map[string]map[int]Action[T]
	err = json.Unmarshal(content, &mapping)
	if err != nil {
		slog.Error("Error unmarshalling maco list: ", "error", err)
	}

	mh.List = mapping

	slog.Debug("Loaded mapping list: ", "count", len(mapping))

	return len(mapping)
}

func (mh *MappingHandler[T]) add(actionType string, devIdentifier string, appIdentifier string, devKey int, val T) {
	if _, exists := mh.List[devIdentifier]; !exists {
		mh.List[devIdentifier] = make(map[string]map[int]Action[T])
	}
	if _, exists := mh.List[devIdentifier][appIdentifier]; !exists {
		mh.List[devIdentifier][appIdentifier] = make(map[int]Action[T])
	}
	mh.List[devIdentifier][appIdentifier][devKey] = Action[T]{Type: actionType, Value: val}
}

// ------------ unused
func (mh *MappingHandler[T]) GetDeviceKeyAction(devIdentifier string, appIdentifier string, devKey int) Action[T] {
	mh.mappingListMutex.Lock()
	defer mh.mappingListMutex.Unlock()

	return mh.List[devIdentifier][appIdentifier][devKey]
}

func (mh *MappingHandler[T]) ClearDeviceKeyAction(devIdentifier string, appIdentifier string, devKey int) {
	mh.mappingListMutex.Lock()
	defer mh.mappingListMutex.Unlock()

	delete(mh.List[devIdentifier][appIdentifier], devKey)
}

func (mh *MappingHandler[T]) ClearAll() {
	mh.mappingListMutex.Lock()
	defer mh.mappingListMutex.Unlock()

	mh.List = make(map[string]map[string]map[int]Action[T])
}

func (mh *MappingHandler[T]) GetMacroForDev(devIdentifier string, appIdentifier string) *int {
	mh.mappingListMutex.Lock()
	defer mh.mappingListMutex.Unlock()
	devKeyList := mh.List[devIdentifier][appIdentifier]
	for devKey, action := range devKeyList {
		if action.Type == "Macro" {
			return &devKey
		}
	}
	return nil
}

func (mh *MappingHandler[T]) GetKeyForDev(devIdentifier string, appIdentifier string, actionValue T) *int {
	mh.mappingListMutex.Lock()
	defer mh.mappingListMutex.Unlock()
	devKeyList := mh.List[devIdentifier][appIdentifier]
	for devKey, act := range devKeyList {
		if act.Type == "Key" {
			return &devKey
		}
	}
	return nil
}

func (mh *MappingHandler[Macro]) GetDeviceAllMacros(devIdentifier string, appIdentifier string) map[int]Macro {
	mh.mappingListMutex.Lock()
	defer mh.mappingListMutex.Unlock()

	macros := make(map[int]Macro)
	for key, action := range mh.List[devIdentifier][appIdentifier] {
		if action.Type == "Macro" {
			macros[key] = action.Value
		}
	}
	return macros
}

func (mh *MappingHandler[Key]) GetDeviceAllKeys(devIdentifier string, appIdentifier string) map[int]Key {
	mh.mappingListMutex.Lock()
	defer mh.mappingListMutex.Unlock()

	macros := make(map[int]Key)
	for key, action := range mh.List[devIdentifier][appIdentifier] {
		if action.Type == "Key" {
			macros[key] = action.Value
		}
	}
	return macros
}

func (mh *MappingHandler[T]) DeviceKeyHasActions(devIdentifier string, appIdentifier string, devKey int) bool {
	mh.mappingListMutex.Lock()
	defer mh.mappingListMutex.Unlock()

	_, exists := mh.List[devIdentifier][appIdentifier][devKey]
	return exists
}

func (mh *MappingHandler[T]) DeviceKeyHasMacro(devIdentifier string, appIdentifier string, devKey int) bool {
	mh.mappingListMutex.Lock()
	defer mh.mappingListMutex.Unlock()

	if action, exists := mh.List[devIdentifier][appIdentifier][devKey]; exists {
		return action.Type == "Macro"
	}
	return false
}

func (mh *MappingHandler[T]) DeviceKeyHasKey(devIdentifier string, appIdentifier string, devKey int) bool {
	mh.mappingListMutex.Lock()
	defer mh.mappingListMutex.Unlock()

	if action, exists := mh.List[devIdentifier][appIdentifier][devKey]; exists {
		return action.Type == "Key"
	}
	return false
}
