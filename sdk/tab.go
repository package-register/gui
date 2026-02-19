package sdk

import (
	"gui/event"

	"github.com/gonutz/wui/v2"
)

// TabContext Tab上下文，暴露给用户回调
type TabContext struct {
	name   string
	panel  *wui.Panel
	app    *App
	events *event.Bus
}

// Name 获取Tab名称
func (t *TabContext) Name() string {
	return t.name
}

// AddLabel 添加标签
func (t *TabContext) AddLabel(text string, x, y, w, h int) *wui.Label {
	label := wui.NewLabel()
	label.SetText(text)
	label.SetBounds(x, y, w, h)
	t.panel.Add(label)
	return label
}

// AddButton 添加按钮
func (t *TabContext) AddButton(text string, x, y, w, h int, onClick func()) *wui.Button {
	btn := wui.NewButton()
	btn.SetText(text)
	btn.SetBounds(x, y, w, h)
	if onClick != nil {
		btn.SetOnClick(onClick)
	}
	t.panel.Add(btn)
	return btn
}

// AddEditLine 添加单行输入框
func (t *TabContext) AddEditLine(x, y, w, h int) *wui.EditLine {
	edit := wui.NewEditLine()
	edit.SetBounds(x, y, w, h)
	t.panel.Add(edit)
	return edit
}

// AddTextEdit 添加多行文本框
func (t *TabContext) AddTextEdit(x, y, w, h int) *wui.TextEdit {
	edit := wui.NewTextEdit()
	edit.SetBounds(x, y, w, h)
	t.panel.Add(edit)
	return edit
}

// AddCheckBox 添加复选框
func (t *TabContext) AddCheckBox(text string, x, y, w, h int, onChange func(bool)) *wui.CheckBox {
	cb := wui.NewCheckBox()
	cb.SetText(text)
	cb.SetBounds(x, y, w, h)
	if onChange != nil {
		cb.SetOnChange(onChange)
	}
	t.panel.Add(cb)
	return cb
}

// AddProgressBar 添加进度条
func (t *TabContext) AddProgressBar(x, y, w, h int) *wui.ProgressBar {
	pb := wui.NewProgressBar()
	pb.SetBounds(x, y, w, h)
	t.panel.Add(pb)
	return pb
}

// AddPanel 添加子面板（用于分组布局）
func (t *TabContext) AddPanel(x, y, w, h int) *wui.Panel {
	p := wui.NewPanel()
	p.SetBounds(x, y, w, h)
	t.panel.Add(p)
	return p
}

// AddSeparator 添加水平分隔线（用Label模拟）
func (t *TabContext) AddSeparator(x, y, w int) *wui.Label {
	sep := wui.NewLabel()
	sep.SetText("────────────────────────────────────────────────────────")
	sep.SetBounds(x, y, w, 2)
	t.panel.Add(sep)
	return sep
}

// Panel 获取底层Panel（高级用法）
func (t *TabContext) Panel() *wui.Panel {
	return t.panel
}

// App 获取应用引用
func (t *TabContext) App() *App {
	return t.app
}

// Events 获取事件总线
func (t *TabContext) Events() *event.Bus {
	return t.events
}

func (t *TabContext) show() {
	t.panel.SetBounds(0, t.app.contentY, t.app.width, t.app.height-t.app.contentY)
}

func (t *TabContext) hide() {
	t.panel.SetBounds(0, t.app.contentY, 0, 0)
}
