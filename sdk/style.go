package sdk

import (
	"github.com/gonutz/wui/v2"
)

// Theme 主题配置
type Theme struct {
	// 颜色
	Background    wui.Color
	Surface       wui.Color
	Foreground    wui.Color
	Primary       wui.Color
	Secondary     wui.Color
	Accent        wui.Color
	Error         wui.Color
	Border        wui.Color

	// 字体
	DefaultFont  string
	HeadingFont  string
	MonoFont     string
	FontSize     int // 负值表示点数

	// 间距
	XSmallPadding int
	SmallPadding  int
	MediumPadding int
	LargePadding  int
	XLargePadding int

	// 边框
	BorderWidth   int
	CornerRadius  int
}

// DefaultTheme 默认现代主题（Material Design 风格）
func DefaultTheme() *Theme {
	return &Theme{
		// 浅色主题配色
		Background:    wui.RGB(250, 250, 250),    // #FAFAFA
		Surface:       wui.RGB(255, 255, 255),    // #FFFFFF
		Foreground:    wui.RGB(33, 33, 33),       // #212121
		Primary:       wui.RGB(103, 80, 164),     // #6750A4 (Material 3)
		Secondary:     wui.RGB(179, 157, 219),    // #B39DDB
		Accent:        wui.RGB(103, 80, 164),      // #6750A4
		Error:         wui.RGB(186, 26, 26),       // #BA1A1A
		Border:        wui.RGB(200, 200, 200),     // #C8C8C8

		// 字体
		DefaultFont:  "微软雅黑",
		HeadingFont:  "微软雅黑",
		MonoFont:     "Consolas",
		FontSize:     -14, // 14pt

		// 间距
		XSmallPadding: 4,
		SmallPadding:  8,
		MediumPadding: 16,
		LargePadding:  24,
		XLargePadding: 32,

		// 边框
		BorderWidth:   1,
		CornerRadius:  4,
	}
}

// DarkTheme 深色主题
func DarkTheme() *Theme {
	return &Theme{
		Background:    wui.RGB(28, 27, 31),       // #1C1B1F
		Surface:       wui.RGB(49, 48, 51),       // #313033
		Foreground:    wui.RGB(230, 225, 229),   // #E6E1E5
		Primary:       wui.RGB(210, 196, 255),    // #D2C4FF
		Secondary:     wui.RGB(137, 124, 176),    // #897CB0
		Accent:        wui.RGB(210, 196, 255),    // #D2C4FF
		Error:         wui.RGB(255, 180, 180),   // #FFB4B4
		Border:        wui.RGB(80, 80, 80),       // #505050

		DefaultFont:  "微软雅黑",
		HeadingFont:  "微软雅黑",
		MonoFont:     "Consolas",
		FontSize:     -14,

		XSmallPadding: 4,
		SmallPadding:  8,
		MediumPadding: 16,
		LargePadding:  24,
		XLargePadding: 32,

		BorderWidth:   1,
		CornerRadius:  4,
	}
}

// ApplyToPanel 将主题应用到面板
func (t *Theme) ApplyToPanel(panel *wui.Panel, borderStyle wui.PanelBorderStyle) {
	panel.SetBorderStyle(borderStyle)
}

// ApplyToButton 将主题应用到按钮
func (t *Theme) ApplyToButton(btn *wui.Button) {
	// 注意：wui 不支持自定义按钮颜色，只能使用默认样式
	// 此函数保留用于未来扩展或自定义绘制
}

// ApplyToLabel 将主题应用到标签
func (t *Theme) ApplyToLabel(label *wui.Label) {
	// 注意：wui 不支持自定义标签颜色
	// 此函数保留用于未来扩展
}

// CreateStyledPanel 创建带主题样式的面板
func (t *Theme) CreateStyledPanel(x, y, w, h int, borderStyle wui.PanelBorderStyle) *wui.Panel {
	panel := wui.NewPanel()
	panel.SetBounds(x, y, w, h)
	panel.SetBorderStyle(borderStyle)
	return panel
}

// CreateStyledLabel 创建带主题样式的标签
func (t *Theme) CreateStyledLabel(text string, x, y, w, h int) *wui.Label {
	label := wui.NewLabel()
	label.SetText(text)
	label.SetBounds(x, y, w, h)
	return label
}

// CreateStyledButton 创建带主题样式的按钮
func (t *Theme) CreateStyledButton(text string, x, y, w, h int, onClick func()) *wui.Button {
	btn := wui.NewButton()
	btn.SetText(text)
	btn.SetBounds(x, y, w, h)
	if onClick != nil {
		btn.SetOnClick(onClick)
	}
	return btn
}

// CreateStyledEditLine 创建带主题样式的输入框
func (t *Theme) CreateStyledEditLine(x, y, w, h int) *wui.EditLine {
	edit := wui.NewEditLine()
	edit.SetBounds(x, y, w, h)
	return edit
}

// CreateStyledTextEdit 创建带主题样式的多行文本框
func (t *Theme) CreateStyledTextEdit(x, y, w, h int, readOnly bool) *wui.TextEdit {
	edit := wui.NewTextEdit()
	edit.SetBounds(x, y, w, h)
	edit.SetReadOnly(readOnly)
	return edit
}

// CreateStyledChatPanel 创建带主题样式的聊天面板
func (t *Theme) CreateStyledChatPanel(x, y, w, h int) *ChatPanel {
	panel := t.CreateStyledPanel(x, y, w, h, wui.PanelBorderSunken)

	padding := t.MediumPadding
	buttonHeight := 32
	inputHeight := 36

	historyHeight := h - inputHeight - buttonHeight - padding*3

	historyEdit := t.CreateStyledTextEdit(padding, padding, w-padding*2, historyHeight, true)
	panel.Add(historyEdit)

	inputY := padding + historyHeight + padding
	inputWidth := w - buttonHeight*2 - padding*3
	inputEdit := t.CreateStyledEditLine(padding, inputY, inputWidth, inputHeight)
	panel.Add(inputEdit)

	btnX := padding + inputWidth + padding
	sendBtn := t.CreateStyledButton("发送", btnX, inputY, buttonHeight*2, inputHeight, nil)
	panel.Add(sendBtn)

	chatPanel := &ChatPanel{
		panel:   panel,
		history: historyEdit,
		input:   inputEdit,
		sendBtn: sendBtn,
	}

	sendBtn.SetOnClick(func() {
		if chatPanel.onSend != nil {
			chatPanel.onSend()
		}
	})

	return chatPanel
}

// WithTheme 配置选项：设置主题
func WithTheme(theme *Theme) Option {
	return func(a *App) {
		a.theme = theme
	}
}

// GetPadding 获取指定级别的内边距
func (t *Theme) GetPadding(level int) int {
	switch level {
	case 0:
		return t.XSmallPadding
	case 1:
		return t.SmallPadding
	case 2:
		return t.MediumPadding
	case 3:
		return t.LargePadding
	case 4:
		return t.XLargePadding
	default:
		return t.MediumPadding
	}
}

// GetSpacing 获取指定级别的间距
func (t *Theme) GetSpacing(level int) int {
	return t.GetPadding(level)
}
