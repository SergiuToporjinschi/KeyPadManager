package gui

import (
	"fmt"
	"main/logger"
	"main/usb"
	"syscall"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/shirou/gopsutil/process"
)

var (
	user32                       = syscall.NewLazyDLL("user32.dll")
	procGetForegroundWindow      = user32.NewProc("GetForegroundWindow")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
)

type RunningApps struct {
	autRefBtn     *widget.Button
	addApp        *widget.Button
	autoRefresher *AutoRefresher
}

func NewRunningApps(_ *GUI) *RunningApps {
	return &RunningApps{}
}

func (r *RunningApps) Build(device *usb.Device) *fyne.Container {
	items := []fyne.CanvasObject{}
	bName := binding.NewString()
	bExe := binding.NewString()
	bPID := binding.NewString()

	err := r.getFocusWindow(bPID, bExe, bName)
	if err != nil {
		logger.Log.Error(err)
	}

	r.autoRefresher = NewAutoRefresher(func() { r.getFocusWindow(bPID, bExe, bName) }, r.updateAutoRefBtn)

	r.autRefBtn = widget.NewButton("Auto refresh", r.autoRefresher.Toggle)
	r.addApp = widget.NewButton("Add app profile", r.addAppProfile)

	items = append(items, r.addApp, r.autRefBtn)
	items = append(items, widget.NewLabel("PID:"), widget.NewLabelWithData(bPID))
	items = append(items, widget.NewLabel("Name:"), widget.NewLabelWithData(bName))
	items = append(items, widget.NewLabel("Exe:"), widget.NewLabelWithData(bExe))

	return container.New(layout.NewFormLayout(), items...)
}
func (r *RunningApps) Destroy() {
	if r.autoRefresher != nil {
		r.autoRefresher.Stop()
	}
}

func (r *RunningApps) updateAutoRefBtn(state bool) {
	if state {
		r.autRefBtn.SetText("Stop auto refresh")
	} else {
		r.autRefBtn.SetText("Auto refresh")
	}
}

func (r *RunningApps) addAppProfile() {
	logger.Log.Info("Create app profile")
}

func (r *RunningApps) getFocusWindow(bPID binding.String, bExe binding.String, bName binding.String) error {
	hwnd := r.getForegroundWindow()
	if hwnd == 0 {
		return fmt.Errorf("no foreground window found")
	}

	pid := r.getWindowThreadProcessId(hwnd)
	bPID.Set(fmt.Sprint(pid))

	// Get detailed information about the foreground process
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return err
	}

	name, err := p.Name()
	if err != nil {
		return err
	}
	bName.Set(name)

	exe, err := p.Exe()
	if err != nil {
		return err
	}
	bExe.Set(exe)
	return nil
}

func (r *RunningApps) getForegroundWindow() syscall.Handle {
	ret, _, _ := procGetForegroundWindow.Call()
	return syscall.Handle(ret)
}

func (r *RunningApps) getWindowThreadProcessId(hwnd syscall.Handle) uint32 {
	var pid uint32
	procGetWindowThreadProcessId.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&pid)))
	return pid
}
