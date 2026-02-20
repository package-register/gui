package sdk

import "github.com/gonutz/wui/v2"

// Layout 布局类型
type Layout int

const (
	LayoutAbsolute Layout = iota // 绝对定位（默认）
	LayoutColumn                // 列布局（垂直堆叠）
	LayoutRow                   // 行布局（水平排列）
	LayoutGrid                  // 网格布局
)

// LayoutConfig 布局配置
type LayoutConfig struct {
	Type      Layout
	Padding   int
	Spacing   int
	Width     int
	Height    int
	GridCols  int // 列数（仅 Grid 布局使用）
}

// LayoutHelper 布局辅助器
type LayoutHelper struct {
	config *LayoutConfig
	panel  *wui.Panel
	x      int
	y      int
	col    int
	row    int
}

// NewLayoutHelper 创建布局辅助器
func NewLayoutHelper(config *LayoutConfig, panel *wui.Panel) *LayoutHelper {
	if config.Padding < 0 {
		config.Padding = 0
	}
	if config.Spacing < 0 {
		config.Spacing = 0
	}

	return &LayoutHelper{
		config: config,
		panel:  panel,
		x:      config.Padding,
		y:      config.Padding,
		col:    0,
		row:    0,
	}
}

// AddChild 添加子组件（根据布局类型自动定位）
func (l *LayoutHelper) AddChild(control wui.Control, w, h int) {
	switch l.config.Type {
	case LayoutRow:
		l.addToRow(control, w, h)
	case LayoutColumn:
		l.addToColumn(control, w, h)
	case LayoutGrid:
		l.addToGrid(control, w, h)
	default: // LayoutAbsolute
		l.addToAbsolute(control, w, h)
	}
}

// AddChildWithPos 指定位置添加子组件（适用于 Grid 或 Absolute 布局）
func (l *LayoutHelper) AddChildWithPos(control wui.Control, x, y, w, h int) {
	control.SetBounds(x, y, w, h)
	l.panel.Add(control)
}

// addToRow 添加到行布局
func (l *LayoutHelper) addToRow(control wui.Control, w, h int) {
	control.SetBounds(l.x, l.y, w, h)
	l.panel.Add(control)
	l.x += w + l.config.Spacing
}

// addToColumn 添加到列布局
func (l *LayoutHelper) addToColumn(control wui.Control, w, h int) {
	control.SetBounds(l.x, l.y, w, h)
	l.panel.Add(control)
	l.y += h + l.config.Spacing
}

// addToGrid 添加到网格布局
func (l *LayoutHelper) addToGrid(control wui.Control, w, h int) {
	if l.config.GridCols <= 0 {
		l.config.GridCols = 2 // 默认 2 列
	}

	// 计算单元格大小
	cellWidth := (l.config.Width - l.config.Padding*2 - l.config.Spacing*(l.config.GridCols-1)) / l.config.GridCols
	cellHeight := h

	// 计算位置
	x := l.config.Padding + l.col*(cellWidth+l.config.Spacing)
	y := l.config.Padding + l.row*(cellHeight+l.config.Spacing)

	control.SetBounds(x, y, cellWidth, cellHeight)
	l.panel.Add(control)

	// 更新行列
	l.col++
	if l.col >= l.config.GridCols {
		l.col = 0
		l.row++
	}
}

// addToAbsolute 添加到绝对定位
func (l *LayoutHelper) addToAbsolute(control wui.Control, w, h int) {
	control.SetBounds(l.x, l.y, w, h)
	l.panel.Add(control)
}

// AddLabel 添加标签（便捷方法）
func (l *LayoutHelper) AddLabel(text string, w, h int) *wui.Label {
	label := wui.NewLabel()
	label.SetText(text)
	l.AddChild(label, w, h)
	return label
}

// AddButton 添加按钮（便捷方法）
func (l *LayoutHelper) AddButton(text string, w, h int, onClick func()) *wui.Button {
	btn := wui.NewButton()
	btn.SetText(text)
	if onClick != nil {
		btn.SetOnClick(onClick)
	}
	l.AddChild(btn, w, h)
	return btn
}

// AddEditLine 添加输入框（便捷方法）
func (l *LayoutHelper) AddEditLine(w, h int) *wui.EditLine {
	edit := wui.NewEditLine()
	l.AddChild(edit, w, h)
	return edit
}

// NewRowLayout 创建行布局辅助器
func NewRowLayout(padding, spacing, width, height int, panel *wui.Panel) *LayoutHelper {
	return NewLayoutHelper(&LayoutConfig{
		Type:    LayoutRow,
		Padding: padding,
		Spacing: spacing,
		Width:   width,
		Height:  height,
	}, panel)
}

// NewColumnLayout 创建列布局辅助器
func NewColumnLayout(padding, spacing, width, height int, panel *wui.Panel) *LayoutHelper {
	return NewLayoutHelper(&LayoutConfig{
		Type:    LayoutColumn,
		Padding: padding,
		Spacing: spacing,
		Width:   width,
		Height:  height,
	}, panel)
}

// NewGridLayout 创建网格布局辅助器
func NewGridLayout(padding, spacing, width, height, cols int, panel *wui.Panel) *LayoutHelper {
	return NewLayoutHelper(&LayoutConfig{
		Type:     LayoutGrid,
		Padding:  padding,
		Spacing:  spacing,
		Width:    width,
		Height:   height,
		GridCols: cols,
	}, panel)
}

// BoxLayout 盒子布局辅助函数
// 自动计算组件位置并添加到面板
func BoxLayout(panel *wui.Panel, padding, spacing int, controls []wui.Control, widths, heights []int, vertical bool) {
	x := padding
	y := padding

	for i, control := range controls {
		w := widths[i]
		h := heights[i]

		control.SetBounds(x, y, w, h)
		panel.Add(control)

		if vertical {
			y += h + spacing
		} else {
			x += w + spacing
		}
	}
}

// GridLayout 网格布局辅助函数
// 将组件排列成指定的列数
func GridLayout(panel *wui.Panel, padding, spacing, width, height, cols int, controls []wui.Control, itemWidths, itemHeights []int) {
	if cols <= 0 {
		cols = 2
	}

	cellWidth := (width - padding*2 - spacing*(cols-1)) / cols
	if len(itemWidths) > 0 {
		cellWidth = itemWidths[0] // 使用第一个组件的宽度
	}

	cellHeight := height
	if len(itemHeights) > 0 {
		cellHeight = itemHeights[0]
	}

	row := 0
	col := 0

	for i, control := range controls {
		w := cellWidth
		h := cellHeight
		if i < len(itemWidths) && itemWidths[i] > 0 {
			w = itemWidths[i]
		}
		if i < len(itemHeights) && itemHeights[i] > 0 {
			h = itemHeights[i]
		}

		x := padding + col*(w+spacing)
		y := padding + row*(h+spacing)

		control.SetBounds(x, y, w, h)
		panel.Add(control)

		col++
		if col >= cols {
			col = 0
			row++
		}
	}
}
