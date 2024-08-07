package gui

import (
	"image/color"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Theme = (*ScsiTheme)(nil)
var ScsiSizeNameText float32 = 12
var ScsiSizeNamePadding float32 = 10
var ScsiSizeNameLineSpacing float32 = 1
var ScsiSizeNameSeparatorThickness float32 = 1
var ScsiMonochrome bool

type ScsiTheme struct {
}

//goland:noinspection GoUnusedExportedFunction
func NewTheme(prefs fyne.Preferences) *ScsiTheme {
	t := &ScsiTheme{}
	ScsiSizeNameText = float32(prefs.FloatWithFallback("sizeNameText", float64(ScsiSizeNameText)))
	ScsiSizeNamePadding = float32(prefs.FloatWithFallback("sizeNamePadding", float64(ScsiSizeNamePadding)))
	ScsiSizeNameLineSpacing = float32(prefs.FloatWithFallback("sizeNameLineSpacing", float64(ScsiSizeNameLineSpacing)))
	ScsiSizeNameSeparatorThickness = float32(prefs.FloatWithFallback("sizeNameSeparatorThickness", float64(ScsiSizeNameSeparatorThickness)))
	ScsiMonochrome = prefs.BoolWithFallback("monochrome", ScsiMonochrome)
	prefs.SetFloat("sizeNameText", float64(ScsiSizeNameText))
	prefs.SetFloat("sizeNamePadding", float64(ScsiSizeNamePadding))
	prefs.SetFloat("sizeNameLineSpacing", float64(ScsiSizeNameLineSpacing))
	prefs.SetFloat("sizeNameSeparatorThickness", float64(ScsiSizeNameSeparatorThickness))
	prefs.SetBool("monochrome", ScsiMonochrome)
	return t
}

//goland:noinspection GoUnusedExportedFunction
func TextPlus() {
	ScsiSizeNameText++
}

//goland:noinspection GoUnusedExportedFunction
func TextMinus() {
	if ScsiSizeNameText > 6 {
		ScsiSizeNameText--
	}
}
func (t ScsiTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	slog.Debug("ScsiTheme.Color", "name", name, "variant", variant)

	switch name {
	case theme.ColorNameSeparator:
		return theme.PrimaryColor()
	case theme.ColorNameScrollBar:
		return theme.PrimaryColor()
	case theme.ColorNameBackground:
		if ScsiMonochrome {
			switch variant {
			case 0: // dark
				return color.Black
			case 1: // light
				return color.White
			}
		}
	case theme.ColorNameForeground:
		if ScsiMonochrome {
			switch variant {
			case 0: // dark
				return color.White
			case 1: // light
				return color.Black
			}
		}
	}
	return theme.DefaultTheme().Color(name, variant)
}
func (t ScsiTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t ScsiTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}
func (t ScsiTheme) Size(name fyne.ThemeSizeName) float32 {

	switch name {
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameInlineIcon:
		return 20
	case theme.SizeNameScrollBar:
		return 5
	case theme.SizeNameScrollBarSmall:
		return 5
	case theme.SizeNameText:
		return 18
	case theme.SizeNameHeadingText:
		return 30
	case theme.SizeNameSubHeadingText:
		return 25
	case theme.SizeNameCaptionText:
		return 15
	case theme.SizeNameInputBorder:
		return 1
	}
	return theme.DefaultTheme().Size(name)
}
