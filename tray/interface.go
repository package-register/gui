package tray

// Adapter 托盘适配器接口（底层实现）
type Adapter interface {
	Initialize(onReady func(), onExit func()) error
	SetIcon(iconBytes []byte)
	SetTitle(title string)
	SetTooltip(tooltip string)
	AddMenuItem(title, tooltip string, handler func()) MenuItem
	AddSeparator()
	Quit()
	IsRunning() bool
}

// MenuItem 菜单项接口
type MenuItem interface {
	SetTitle(title string)
	SetTooltip(tooltip string)
	SetIcon(iconBytes []byte)
	Check()
	Uncheck()
	IsChecked() bool
	Disable()
	Enable()
	IsEnabled() bool
	OnClick(handler func())
}

// Tray 托盘控制器（暴露给用户的接口）
type Tray struct {
	adapter Adapter
	icon    []byte
	tooltip string
	running bool
	pending []func() // 待执行的操作队列
}

// NewTray 创建托盘控制器
func NewTray() *Tray {
	return &Tray{
		adapter: NewFyneAdapter(),
	}
}

// SetIcon 设置图标
func (t *Tray) SetIcon(icon []byte) {
	t.icon = icon
	if t.running {
		t.adapter.SetIcon(icon)
	}
}

// SetTooltip 设置提示
func (t *Tray) SetTooltip(tooltip string) {
	t.tooltip = tooltip
	if t.running {
		t.adapter.SetTooltip(tooltip)
	}
}

// AddMenuItem 添加菜单项
func (t *Tray) AddMenuItem(title, tooltip string, handler func()) MenuItem {
	if t.running {
		return t.adapter.AddMenuItem(title, tooltip, handler)
	}
	// 未就绪，缓存操作
	t.pending = append(t.pending, func() {
		t.adapter.AddMenuItem(title, tooltip, handler)
	})
	return nil
}

// AddSeparator 添加分隔符
func (t *Tray) AddSeparator() {
	if t.running {
		t.adapter.AddSeparator()
		return
	}
	t.pending = append(t.pending, func() {
		t.adapter.AddSeparator()
	})
}

// Quit 退出
func (t *Tray) Quit() {
	if t.running {
		t.adapter.Quit()
		t.running = false
	}
}

// Start 启动托盘
func (t *Tray) Start() error {
	return t.adapter.Initialize(func() {
		t.running = true
		// 设置图标和提示
		if len(t.icon) > 0 {
			t.adapter.SetIcon(t.icon)
		}
		if t.tooltip != "" {
			t.adapter.SetTooltip(t.tooltip)
		}
		// 执行缓存的操作
		for _, fn := range t.pending {
			fn()
		}
		t.pending = nil
	}, func() {
		t.running = false
	})
}
