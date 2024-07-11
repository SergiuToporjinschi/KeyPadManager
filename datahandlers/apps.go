package datahandlers

import (
	"encoding/json"
	"log/slog"
	"os"
	"sync"
)

type Application struct {
	Name    string
	ExePath string
}

var onceAppHandler sync.Once
var instanceAppHandler *AppsHandler

type AppsHandler struct {
	appListMutex sync.Mutex
	appList      []Application
}

func GetInstance() *AppsHandler {
	onceAppHandler.Do(func() {
		instanceAppHandler = &AppsHandler{}
	})
	return instanceAppHandler
}

func (ah *AppsHandler) Add(app Application) bool {
	slog.Debug("Adding app: ", "app", app)

	ah.appListMutex.Lock()
	defer ah.appListMutex.Unlock()

	return ah.add(app)
}

func (ah *AppsHandler) AddAll(apps []Application) int {
	slog.Debug("Adding app list: ", "apps", apps)

	ah.appListMutex.Lock()
	defer ah.appListMutex.Unlock()

	return ah.addAll(apps)
}

func (ah *AppsHandler) Remove(app Application) bool {
	slog.Debug("Removing app: ", "app", app)

	ah.appListMutex.Lock()
	defer ah.appListMutex.Unlock()

	for i, a := range ah.appList {
		if a.ExePath == app.ExePath && a.Name == app.Name {
			ah.appList = append(ah.appList[:i], ah.appList[i+1:]...)
			return true
		}
	}
	return false
}

func (ah *AppsHandler) RemoveByExePath(exePath string) bool {
	slog.Debug("Removing app by exe path: ", "exePath", exePath)

	ah.appListMutex.Lock()
	defer ah.appListMutex.Unlock()

	for i, a := range ah.appList {
		if a.ExePath == exePath {
			ah.appList = append(ah.appList[:i], ah.appList[i+1:]...)
			return true
		}
	}
	return false
}

func (ah *AppsHandler) Clear() {
	slog.Debug("Clearing app list")

	ah.appListMutex.Lock()
	defer ah.appListMutex.Unlock()

	ah.appList = []Application{}
}

func (ah *AppsHandler) LoadAppList() int {
	ah.Clear()

	ah.appListMutex.Lock()
	defer ah.appListMutex.Unlock()

	content, err := os.ReadFile("apps.json")
	if os.IsNotExist(err) {
		slog.Warn("File not found: ", "error", err)
	} else if err != nil {
		slog.Error("Error opening file: ", "error", err)
	}

	var apps []Application
	err = json.Unmarshal(content, &apps)
	if err != nil {
		slog.Error("Error unmarshalling app list: ", "error", err)
	}

	ah.addAll(apps)

	slog.Debug("Loaded app list: ", "count", len(apps))

	return len(apps)
}

func (ah *AppsHandler) SaveAppList() {
	slog.Debug("Saving app list started")

	ah.appListMutex.Lock()
	defer ah.appListMutex.Unlock()

	jsonApps, err := json.MarshalIndent(ah.appList, "", "  ")
	if err != nil {
		slog.Error("Error marshalling app list: ", "error", err)
		return
	}

	file, err := os.OpenFile("apps.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
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

	slog.Debug("Saving app list finished", "count", len(ah.appList))
}

func (ah *AppsHandler) GetAppList() []Application {
	ah.appListMutex.Lock()
	defer ah.appListMutex.Unlock()

	return ah.appList
}

func (ah *AppsHandler) GetAppByExePath(exePath string) *Application {
	ah.appListMutex.Lock()
	defer ah.appListMutex.Unlock()

	for _, a := range ah.appList {
		if a.ExePath == exePath {
			return &a
		}
	}
	return nil
}

func (ah *AppsHandler) GetByIndex(index int) *Application {
	ah.appListMutex.Lock()
	defer ah.appListMutex.Unlock()

	if index >= 0 && index < len(ah.appList) {
		return &ah.appList[index]
	}
	return nil
}

func (ah *AppsHandler) RemoveByIndex(index int) bool {
	ah.appListMutex.Lock()
	defer ah.appListMutex.Unlock()

	if index >= 0 && index < len(ah.appList) {
		ah.appList = append(ah.appList[:index], ah.appList[index+1:]...)
		return true
	}
	return false
}

func (ah *AppsHandler) Size() int {
	ah.appListMutex.Lock()
	defer ah.appListMutex.Unlock()

	return len(ah.appList)
}

func (ah *AppsHandler) IsEmpty() bool {
	ah.appListMutex.Lock()
	defer ah.appListMutex.Unlock()

	return len(ah.appList) == 0
}

func (ah *AppsHandler) GetExePaths() []string {
	ah.appListMutex.Lock()
	defer ah.appListMutex.Unlock()

	exePaths := make([]string, len(ah.appList))
	for i, a := range ah.appList {
		exePaths[i] = a.ExePath
	}
	return exePaths
}

func (ah *AppsHandler) containsExePath(exePath string) bool {
	for _, a := range ah.appList {
		if a.ExePath == exePath {
			return true
		}
	}
	return false
}

func (ah *AppsHandler) ContainsExePath(exePath string) bool {
	ah.appListMutex.Lock()
	defer ah.appListMutex.Unlock()

	return ah.containsExePath(exePath)
}

func (ah *AppsHandler) add(app Application) bool {
	if app.ExePath == "" {
		slog.Warn("Empty exe path")
		return false
	}
	if !ah.containsExePath(app.ExePath) {
		ah.appList = append(ah.appList, app)
		return true
	}
	return false
}
func (ah *AppsHandler) addAll(apps []Application) int {
	addedCnt := 0
	for _, app := range apps {
		if ah.add(app) {
			addedCnt++
		}
	}
	return addedCnt
}
