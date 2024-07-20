package datahandlers

import (
	"encoding/json"
	"log/slog"
	"os"
	"sync"

	"github.com/google/uuid"
)

const (
	macroFileName = "macro.json"
)

type Macro struct {
	ID   *uuid.UUID `json:"id" validate:"required"`
	Name string     `json:"name" validate:"required"`
}

var onceMacroHandler sync.Once
var instanceMacroHandler *MacrosHandler

type MacrosHandler struct {
	macroListMutex sync.Mutex
	macroList      []Macro
}

func GetMacrosHandlerInstance() *MacrosHandler {
	onceMacroHandler.Do(func() {
		instanceMacroHandler = &MacrosHandler{}
	})

	return instanceMacroHandler
}

func (mh *MacrosHandler) Add(macro Macro) bool {
	slog.Debug("Adding macro: ", "macro", macro)

	mh.macroListMutex.Lock()
	defer mh.macroListMutex.Unlock()

	return mh.add(macro)
}

func (mh *MacrosHandler) AddAll(macros []Macro) int {
	slog.Debug("Adding macro list: ", "macros", macros)

	mh.macroListMutex.Lock()
	defer mh.macroListMutex.Unlock()

	return mh.addAll(macros)
}

func (mh *MacrosHandler) Remove(macro Macro) bool {
	slog.Debug("Removing macro: ", "macro", macro)

	mh.macroListMutex.Lock()
	defer mh.macroListMutex.Unlock()

	for i, a := range mh.macroList {
		if a.Name == macro.Name {
			mh.macroList = append(mh.macroList[:i], mh.macroList[i+1:]...)
			return true
		}
	}
	return false
}

func (mh *MacrosHandler) Clear() {
	slog.Debug("Clearing macro list")

	mh.macroListMutex.Lock()
	defer mh.macroListMutex.Unlock()

	mh.macroList = []Macro{}
}

func (mh *MacrosHandler) LoadMacroList() int {
	mh.Clear()

	mh.macroListMutex.Lock()
	defer mh.macroListMutex.Unlock()

	content, err := os.ReadFile(macroFileName)
	if os.IsNotExist(err) {
		slog.Warn("File not found: ", "error", err)
		return 0
	} else if err != nil {
		slog.Error("Error opening file: ", "error", err)
	}

	var macros []Macro
	err = json.Unmarshal(content, &macros)
	if err != nil {
		slog.Error("Error unmarshalling maco list: ", "error", err)
	}

	mh.addAll(macros)

	slog.Debug("Loaded macro list: ", "count", len(macros))

	return len(macros)
}

func (mh *MacrosHandler) SaveMacroList() {
	slog.Debug("Saving macro list started")

	mh.macroListMutex.Lock()
	defer mh.macroListMutex.Unlock()

	jsonApps, err := json.MarshalIndent(mh.macroList, "", "  ")
	if err != nil {
		slog.Error("Error marshalling app list: ", "error", err)
		return
	}

	file, err := os.OpenFile(macroFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		slog.Error("Error opening file: ", "error", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(string(jsonApps))
	if err != nil {
		slog.Error("Error writing to file: ", "error", err)
		return
	}

	slog.Debug("Saving macro list finished", "count", len(mh.macroList))
}

func (mh *MacrosHandler) GetMacroList() []Macro {
	mh.macroListMutex.Lock()
	defer mh.macroListMutex.Unlock()

	return mh.macroList
}

func (mh *MacrosHandler) GetByID(id *uuid.UUID) *Macro {
	mh.macroListMutex.Lock()
	defer mh.macroListMutex.Unlock()

	if id == nil {
		slog.Warn("ID is nil")
		return nil
	}
	for _, a := range mh.macroList {
		if a.ID == id {
			return &a
		}
	}
	return nil
}

func (mh *MacrosHandler) GetByIndex(index int) *Macro {
	mh.macroListMutex.Lock()
	defer mh.macroListMutex.Unlock()

	if index >= 0 && index < len(mh.macroList) {
		return &mh.macroList[index]
	}
	return nil
}

func (mh *MacrosHandler) RemoveByIndex(index int) bool {
	mh.macroListMutex.Lock()
	defer mh.macroListMutex.Unlock()

	if index >= 0 && index < len(mh.macroList) {
		mh.macroList = append(mh.macroList[:index], mh.macroList[index+1:]...)
		return true
	}
	return false
}

func (mh *MacrosHandler) Size() int {
	mh.macroListMutex.Lock()
	defer mh.macroListMutex.Unlock()

	return len(mh.macroList)
}

func (mh *MacrosHandler) IsEmpty() bool {
	mh.macroListMutex.Lock()
	defer mh.macroListMutex.Unlock()

	return len(mh.macroList) == 0
}

func (mh *MacrosHandler) ContainsID(id *uuid.UUID) bool {
	if id == nil {
		slog.Warn("ID is nil")
		return false
	}
	mh.macroListMutex.Lock()
	defer mh.macroListMutex.Unlock()

	return mh.containsID(id)
}

func (mh *MacrosHandler) containsID(uuid *uuid.UUID) bool {
	if uuid == nil {
		slog.Warn("ID is nil")
		return false
	}
	for _, a := range mh.macroList {
		if a.ID == uuid {
			return true
		}
	}
	return false
}

func (mh *MacrosHandler) add(macro Macro) bool {
	if macro.ID == nil {
		*macro.ID = uuid.New()
	}

	if !mh.containsID(macro.ID) {
		mh.macroList = append(mh.macroList, macro)
		return true
	}
	return false
}

func (mh *MacrosHandler) addAll(macros []Macro) int {
	addedCnt := 0
	for _, app := range macros {
		if mh.add(app) {
			addedCnt++
		}
	}
	return addedCnt
}
