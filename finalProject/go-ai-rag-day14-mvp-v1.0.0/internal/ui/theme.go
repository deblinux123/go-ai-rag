package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type appTheme struct {
	variant fyne.ThemeVariant
}

var _ fyne.Theme = (*appTheme)(nil)

func newAppTheme(variant fyne.ThemeVariant) fyne.Theme {
	return &appTheme{variant: variant}
}

func (t *appTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	_ = variant

	if name == theme.ColorNamePrimary {
		return color.NRGBA{R: 124, G: 92, B: 252, A: 255}
	}

	if name == theme.ColorNameSelection {
		return color.NRGBA{R: 124, G: 92, B: 252, A: 72}
	}

	return theme.DefaultTheme().Color(name, t.variant)
}

func (t *appTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *appTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *appTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
