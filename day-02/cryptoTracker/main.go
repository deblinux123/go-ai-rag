package main

import (
	"context"
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ---------- API types & client ----------

type CryptoSymbol struct {
	Symbol                string `json:"symbol"`
	Last                  string `json:"last"`
	Lowest                string `json:"lowest"`
	Highest               string `json:"highest"`
	Date                  string `json:"date"`
	DailyChangePercentage string `json:"daily_change_percentage"`
	SourceExchange        string `json:"source_exchange"`
}

type CryptoResponse struct {
	Status  string         `json:"status"`
	Symbols []CryptoSymbol `json:"symbols"`
}

type CryptoClient struct {
	apiKey     string
	httpClient *http.Client
}

func NewCryptoClient(apiKey string) *CryptoClient {
	return &CryptoClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *CryptoClient) GetSymbol(ctx context.Context, symbol string) (*CryptoSymbol, error) {
	url := "https://api.freecryptoapi.com/v1/getData?symbol=" + symbol

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %s: %s", resp.Status, string(body))
	}

	var data CryptoResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	if len(data.Symbols) == 0 {
		return nil, fmt.Errorf("no data returned for %s", symbol)
	}

	return &data.Symbols[0], nil
}

// ---------- Custom theme ----------

// appTheme wraps Fyne's built-in theme, fixes a chosen light/dark variant
// (independent of the OS setting), and swaps in a custom accent color.
// Toggling calls app.Settings().SetTheme again to force a full repaint.
type appTheme struct {
	variant fyne.ThemeVariant
}

var (
	accent     = color.NRGBA{R: 0x7C, G: 0x5C, B: 0xFC, A: 0xFF} // violet
	accentSoft = color.NRGBA{R: 0x7C, G: 0x5C, B: 0xFC, A: 0x22}
	upColor    = color.NRGBA{R: 0x22, G: 0xC5, B: 0x5E, A: 0xFF}
	downColor  = color.NRGBA{R: 0xF4, G: 0x43, B: 0x36, A: 0xFF}
	darkBg     = color.NRGBA{R: 0x14, G: 0x14, B: 0x1C, A: 0xFF}
	darkPanel  = color.NRGBA{R: 0x1E, G: 0x1E, B: 0x2A, A: 0xFF}
	lightBg    = color.NRGBA{R: 0xF6, G: 0xF6, B: 0xFB, A: 0xFF}
	lightPanel = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
)

func (t *appTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	v := t.variant
	switch name {
	case theme.ColorNamePrimary:
		return accent
	case theme.ColorNameSelection:
		return accentSoft
	case theme.ColorNameBackground:
		if v == theme.VariantDark {
			return darkBg
		}
		return lightBg
	case theme.ColorNameButton, theme.ColorNameInputBackground:
		if v == theme.VariantDark {
			return darkPanel
		}
		return lightPanel
	case theme.ColorNameSuccess:
		return upColor
	case theme.ColorNameError:
		return downColor
	}
	return theme.DefaultTheme().Color(name, v)
}

func (t *appTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *appTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *appTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 10
	case theme.SizeNameInlineIcon:
		return 22
	case theme.SizeNameHeadingText:
		return 26
	}
	return theme.DefaultTheme().Size(name)
}

// ---------- UI ----------

type trackerUI struct {
	app    fyne.App
	client *CryptoClient
	theme  *appTheme

	statusDot   *canvas.Circle
	statusLabel *widget.Label

	priceValue  *canvas.Text
	changeValue *canvas.Text
	lowValue    *widget.Label
	highValue   *widget.Label
	exchgValue  *widget.Label
	updated     *widget.Label

	progressBar  *widget.ProgressBarInfinite
	cryptoSelect *widget.Select
	refreshBtn   *widget.Button
	themeBtn     *widget.Button
}

func newTrackerUI(a fyne.App, client *CryptoClient, th *appTheme) *trackerUI {
	t := &trackerUI{app: a, client: client, theme: th}

	t.statusDot = canvas.NewCircle(theme.Color(theme.ColorNameDisabled))
	t.statusLabel = widget.NewLabel("Ready")

	t.priceValue = canvas.NewText("--", accent)
	t.priceValue.TextSize = 34
	t.priceValue.TextStyle = fyne.TextStyle{Bold: true}

	t.changeValue = canvas.NewText("--", theme.Color(theme.ColorNameForeground))
	t.changeValue.TextSize = 18
	t.changeValue.TextStyle = fyne.TextStyle{Bold: true}

	t.lowValue = widget.NewLabel("--")
	t.highValue = widget.NewLabel("--")
	t.exchgValue = widget.NewLabel("--")
	t.updated = widget.NewLabel("")
	t.updated.TextStyle = fyne.TextStyle{Italic: true}

	t.progressBar = widget.NewProgressBarInfinite()
	t.progressBar.Hide()

	t.cryptoSelect = widget.NewSelect([]string{"BTC", "ETH", "BNB", "SOL"}, nil)
	t.cryptoSelect.SetSelected("BTC")

	t.refreshBtn = widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), t.onRefresh)
	t.refreshBtn.Importance = widget.HighImportance

	t.themeBtn = widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), t.toggleTheme)

	return t
}

func (t *trackerUI) toggleTheme() {
	if t.theme.variant == theme.VariantDark {
		t.theme.variant = theme.VariantLight
	} else {
		t.theme.variant = theme.VariantDark
	}
	// Re-applying the theme forces Fyne to repaint every widget with
	// the new colors.
	t.app.Settings().SetTheme(t.theme)
}

func (t *trackerUI) setStatus(text string, dotColor color.Color) {
	t.statusLabel.SetText(text)
	t.statusDot.FillColor = dotColor
	t.statusDot.Refresh()
}

// onRefresh runs on the UI goroutine (it's a widget callback), but it
// immediately hands the slow work off to a background goroutine so the
// UI thread stays free to actually paint the progress bar.
func (t *trackerUI) onRefresh() {
	symbol := t.cryptoSelect.Selected

	t.progressBar.Show()
	t.setStatus("Loading "+symbol+"...", theme.Color(theme.ColorNameWarning))
	t.refreshBtn.Disable()
	t.cryptoSelect.Disable()

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		price, err := t.client.GetSymbol(ctx, symbol)

		// Hop back onto the UI goroutine to touch widgets safely.
		fyne.Do(func() {
			t.progressBar.Hide()
			t.refreshBtn.Enable()
			t.cryptoSelect.Enable()

			if err != nil {
				t.setStatus("Error: "+err.Error(), downColor)
				return
			}

			t.setStatus("Updated", upColor)
			t.priceValue.Text = "$" + price.Last
			t.priceValue.Refresh()

			pct, _ := strconv.ParseFloat(strings.TrimPrefix(price.DailyChangePercentage, "+"), 64)
			arrow := "▲"
			col := upColor
			if pct < 0 {
				arrow = "▼"
				col = downColor
			}
			t.changeValue.Text = arrow + " " + price.DailyChangePercentage + "%"
			t.changeValue.Color = col
			t.changeValue.Refresh()

			t.lowValue.SetText("$" + price.Lowest)
			t.highValue.SetText("$" + price.Highest)
			t.exchgValue.SetText(price.SourceExchange)
			t.updated.SetText("Last updated " + time.Now().Format("15:04:05"))
		})
	}()
}

func tile(title string, value fyne.CanvasObject) *widget.Card {
	titleLabel := widget.NewLabel(title)
	titleLabel.Importance = widget.LowImportance
	return widget.NewCard("", "", container.NewVBox(titleLabel, value))
}

func (t *trackerUI) buildContent() fyne.CanvasObject {
	titleLabel := widget.NewLabelWithStyle("₿  Crypto Tracker", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	header := container.NewBorder(nil, nil, titleLabel, t.themeBtn)

	controls := container.NewBorder(nil, nil, widget.NewLabel("Symbol"), t.refreshBtn, t.cryptoSelect)

	statusRow := container.NewHBox(
		container.New(&fixedSizeLayout{size: fyne.NewSize(10, 10)}, t.statusDot),
		t.statusLabel,
	)

	heroCard := widget.NewCard("", "", container.NewVBox(
		t.priceValue,
		t.changeValue,
	))

	statsGrid := container.NewGridWithColumns(2,
		tile("Low", t.lowValue),
		tile("High", t.highValue),
	)

	exchangeCard := widget.NewCard("", "", container.NewBorder(nil, nil,
		widget.NewLabel("Exchange"), nil, t.exchgValue))

	return container.NewPadded(container.NewVBox(
		header,
		widget.NewSeparator(),
		controls,
		t.progressBar,
		statusRow,
		widget.NewSeparator(),
		heroCard,
		statsGrid,
		exchangeCard,
		t.updated,
	))
}

// fixedSizeLayout pins a single child to an exact size — used to keep the
// status "dot" small and circular regardless of container growth.
type fixedSizeLayout struct {
	size fyne.Size
}

func (f *fixedSizeLayout) MinSize(_ []fyne.CanvasObject) fyne.Size { return f.size }

func (f *fixedSizeLayout) Layout(objs []fyne.CanvasObject, _ fyne.Size) {
	for _, o := range objs {
		o.Resize(f.size)
	}
}

// ---------- main ----------

func main() {
	apiKey := os.Getenv("FREECRYPTO_API_KEY")

	a := app.New()

	th := &appTheme{variant: theme.VariantDark}
	a.Settings().SetTheme(th)

	w := a.NewWindow("Crypto Price Tracker")

	ui := newTrackerUI(a, NewCryptoClient(apiKey), th)
	w.SetContent(ui.buildContent())
	w.Resize(fyne.NewSize(440, 560))

	if apiKey == "" {
		ui.setStatus("Warning: FREECRYPTO_API_KEY not set", downColor)
	}

	w.ShowAndRun()
}
