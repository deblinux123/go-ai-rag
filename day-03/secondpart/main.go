package main

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ---------- Custom theme ----------

// appTheme wraps Fyne's built-in theme, fixes a chosen light/dark variant
// (independent of the OS setting), and applies a custom accent color.
type appTheme struct {
	variant fyne.ThemeVariant
}

var (
	accent     = color.NRGBA{R: 0x6C, G: 0x5C, B: 0xE7, A: 0xFF} // indigo
	accentSoft = color.NRGBA{R: 0x6C, G: 0x5C, B: 0xE7, A: 0x24}
	userBubble = color.NRGBA{R: 0x6C, G: 0x5C, B: 0xE7, A: 0xFF}
	aiDark     = color.NRGBA{R: 0x24, G: 0x24, B: 0x32, A: 0xFF}
	aiLight    = color.NRGBA{R: 0xEC, G: 0xEC, B: 0xF3, A: 0xFF}
	darkBg     = color.NRGBA{R: 0x12, G: 0x12, B: 0x1A, A: 0xFF}
	darkPanel  = color.NRGBA{R: 0x1B, G: 0x1B, B: 0x27, A: 0xFF}
	lightBg    = color.NRGBA{R: 0xF6, G: 0xF6, B: 0xFB, A: 0xFF}
	lightPanel = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	onlineDot  = color.NRGBA{R: 0x22, G: 0xC5, B: 0x5E, A: 0xFF}
	busyDot    = color.NRGBA{R: 0xF5, G: 0xA6, B: 0x23, A: 0xFF}
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
		return 20
	}
	return theme.DefaultTheme().Size(name)
}

// ---------- Chat bubble ----------

// bubbleMaxWidthLayout caps a single child's width so long messages wrap
// instead of stretching edge to edge like a text box.
type bubbleMaxWidthLayout struct {
	maxWidth float32
}

func (b *bubbleMaxWidthLayout) MinSize(objs []fyne.CanvasObject) fyne.Size {
	if len(objs) == 0 {
		return fyne.NewSize(0, 0)
	}
	ms := objs[0].MinSize()
	if ms.Width > b.maxWidth {
		ms.Width = b.maxWidth
	}
	return ms
}

func (b *bubbleMaxWidthLayout) Layout(objs []fyne.CanvasObject, size fyne.Size) {
	for _, o := range objs {
		w := b.maxWidth
		if size.Width < w {
			w = size.Width
		}
		o.Resize(fyne.NewSize(w, o.MinSize().Height))
	}
}

func newBubble(text string, fromUser bool, th *appTheme) fyne.CanvasObject {
	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextWrapWord

	bg := canvas.NewRectangle(aiDark)
	bg.CornerRadius = 14
	textColor := theme.Color(theme.ColorNameForeground)

	if fromUser {
		bg.FillColor = userBubble
		textColor = color.White
	} else if th.variant == theme.VariantDark {
		bg.FillColor = aiDark
	} else {
		bg.FillColor = aiLight
	}

	rich := widget.NewLabel(text)
	rich.Wrapping = fyne.TextWrapWord
	_ = textColor // Label color follows theme automatically; kept for clarity.
	_ = label

	padded := container.NewPadded(rich)
	sized := container.New(&bubbleMaxWidthLayout{maxWidth: 460}, container.NewStack(bg, padded))

	if fromUser {
		return container.NewHBox(layout.NewSpacer(), sized)
	}
	return container.NewHBox(sized, layout.NewSpacer())
}

// ---------- Main UI ----------

type assistantUI struct {
	app   fyne.App
	theme *appTheme

	messages   *fyne.Container
	scroll     *container.Scroll
	input      *widget.Entry
	sendBtn    *widget.Button
	statusDot  *canvas.Circle
	statusText *widget.Label
	modelSel   *widget.Select
}

func newAssistantUI(a fyne.App, th *appTheme) *assistantUI {
	u := &assistantUI{app: a, theme: th}

	u.messages = container.NewVBox()
	u.scroll = container.NewVScroll(u.messages)

	u.input = widget.NewMultiLineEntry()
	u.input.SetPlaceHolder("Message the assistant...")
	u.input.Wrapping = fyne.TextWrapWord

	u.sendBtn = widget.NewButtonWithIcon("Send", theme.MailSendIcon(), u.onSend)
	u.sendBtn.Importance = widget.HighImportance

	u.statusDot = canvas.NewCircle(onlineDot)
	u.statusText = widget.NewLabel("Ready")

	u.modelSel = widget.NewSelect([]string{"llama3", "mistral", "phi3"}, nil)
	u.modelSel.SetSelected("llama3")

	u.appendBubble("Hi! I'm your local AI assistant. Ask me anything, or try a quick action below. 👋", false)

	return u
}

func (u *assistantUI) appendBubble(text string, fromUser bool) {
	u.messages.Add(newBubble(text, fromUser, u.theme))
	u.messages.Refresh()
	u.scroll.ScrollToBottom()
}

func (u *assistantUI) setStatus(text string, dotColor color.Color) {
	u.statusText.SetText(text)
	u.statusDot.FillColor = dotColor
	u.statusDot.Refresh()
}

func (u *assistantUI) toggleTheme() {
	if u.theme.variant == theme.VariantDark {
		u.theme.variant = theme.VariantLight
	} else {
		u.theme.variant = theme.VariantDark
	}
	u.app.Settings().SetTheme(u.theme)
}

func (u *assistantUI) onSend() {
	text := u.input.Text
	if text == "" {
		return
	}

	u.appendBubble(text, true)
	u.input.SetText("")
	u.sendBtn.Disable()
	u.setStatus("Thinking...", busyDot)

	typingIdx := len(u.messages.Objects)
	u.appendBubble("● ● ●", false)

	go func() {
		// TODO: replace with a real call to Ollama / your model backend.
		time.Sleep(2 * time.Second)

		fyne.Do(func() {
			u.messages.Remove(u.messages.Objects[typingIdx])
			u.appendBubble(
				"Hello! 👋\n\nThis is a placeholder response.\n\nNext step: wire this up to Ollama.",
				false,
			)
			u.sendBtn.Enable()
			u.setStatus("Ready", onlineDot)
		})
	}()
}

func (u *assistantUI) quickAction(prompt string) {
	u.input.SetText(prompt)
	u.input.FocusGained()
}

func chip(label string, onTap func()) *widget.Button {
	b := widget.NewButton(label, onTap)
	b.Importance = widget.LowImportance
	return b
}

func (u *assistantUI) buildContent() fyne.CanvasObject {
	// ---- Header ----
	title := widget.NewLabelWithStyle("🤖  AI Desktop Assistant", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	themeBtn := widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), u.toggleTheme)
	settingsBtn := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {})
	header := container.NewBorder(nil, nil, title, container.NewHBox(u.modelSel, themeBtn, settingsBtn))

	// ---- Sidebar ----
	newChatBtn := widget.NewButtonWithIcon("New Chat", theme.ContentAddIcon(), func() {
		u.messages.Objects = nil
		u.appendBubble("New conversation started. What can I help with?", false)
	})
	newChatBtn.Importance = widget.HighImportance

	historyBtn := widget.NewButtonWithIcon("Chat History", theme.HistoryIcon(), func() {})
	modelsBtn := widget.NewButtonWithIcon("Models", theme.ComputerIcon(), func() {})

	sidebar := container.NewVBox(
		newChatBtn,
		widget.NewSeparator(),
		historyBtn,
		modelsBtn,
	)
	sidebarWrap := container.NewBorder(container.NewPadded(sidebar), nil, nil, nil)

	// ---- Quick actions ----
	quickActions := container.NewHBox(
		chip("💡 Explain code", func() { u.quickAction("Explain this code:\n\n") }),
		chip("📝 Summarize", func() { u.quickAction("Summarize the following:\n\n") }),
		chip("✉️ Write email", func() { u.quickAction("Write a professional email about: ") }),
		chip("🧹 Clear chat", func() {
			u.messages.Objects = nil
			u.appendBubble("Chat cleared. Ready when you are.", false)
		}),
	)

	// ---- Chat area ----
	chatArea := container.NewBorder(quickActions, nil, nil, nil, u.scroll)

	// ---- Input bar ----
	u.input.SetMinRowsVisible(2)
	inputBar := container.NewBorder(nil, nil, nil, u.sendBtn, u.input)
	inputCard := container.NewPadded(inputBar)

	mainContent := container.NewBorder(nil, inputCard, nil, nil, chatArea)

	// ---- Footer ----
	statusRow := container.NewHBox(
		container.New(&bubbleMaxWidthLayout{maxWidth: 12}, u.statusDot),
		u.statusText,
	)
	footer := container.NewHBox(statusRow, layout.NewSpacer(), widget.NewLabel("v0.1 · local"))

	return container.NewBorder(
		container.NewPadded(header),
		container.NewPadded(footer),
		sidebarWrap,
		nil,
		container.NewPadded(mainContent),
	)
}

func main() {
	a := app.New()

	th := &appTheme{variant: theme.VariantDark}
	a.Settings().SetTheme(th)

	w := a.NewWindow("AI Desktop Assistant")
	w.Resize(fyne.NewSize(860, 620))

	ui := newAssistantUI(a, th)
	w.SetContent(ui.buildContent())

	w.ShowAndRun()
}
