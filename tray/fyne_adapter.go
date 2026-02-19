package tray

import (
	"sync"

	"fyne.io/systray"
)

// FyneAdapter fyne.io/systray 的适配器实现
type FyneAdapter struct {
	running bool
	mutex   sync.RWMutex

	// 存储菜单项引用
	menuItems []MenuItem
}

// NewFyneAdapter 创建新的Fyne适配器
func NewFyneAdapter() *FyneAdapter {
	return &FyneAdapter{
		menuItems: make([]MenuItem, 0),
	}
}

// Initialize 初始化托盘
func (f *FyneAdapter) Initialize(onReady func(), onExit func()) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.running {
		return nil // 已经初始化
	}

	// 使用RunWithExternalLoop避免阻塞主线程
	systray.RunWithExternalLoop(onReady, onExit)
	f.running = true

	return nil
}

// SetIcon 设置托盘图标
func (f *FyneAdapter) SetIcon(iconBytes []byte) {
	systray.SetIcon(iconBytes)
}

// SetTitle 设置托盘标题（Windows不支持，但提供接口兼容性）
func (f *FyneAdapter) SetTitle(title string) {
	// Windows不支持托盘标题，但为了接口兼容性保留此方法
	// 在macOS和Linux上有效
}

// SetTooltip 设置鼠标悬停提示
func (f *FyneAdapter) SetTooltip(tooltip string) {
	systray.SetTooltip(tooltip)
}

// AddMenuItem 添加菜单项
func (f *FyneAdapter) AddMenuItem(title, tooltip string, handler func()) MenuItem {
	item := &FyneMenuItem{
		item:    systray.AddMenuItem(title, tooltip),
		handler: handler,
	}

	// 启动goroutine监听点击事件
	if handler != nil {
		go func() {
			for range item.item.ClickedCh {
				handler()
			}
		}()
	}

	f.mutex.Lock()
	f.menuItems = append(f.menuItems, item)
	f.mutex.Unlock()

	return item
}

// AddSeparator 添加分隔符
func (f *FyneAdapter) AddSeparator() {
	systray.AddSeparator()
}

// Quit 退出托盘
func (f *FyneAdapter) Quit() {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.running {
		systray.Quit()
		f.running = false
	}
}

// IsRunning 检查托盘是否正在运行
func (f *FyneAdapter) IsRunning() bool {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.running
}

// FyneMenuItem fyne菜单项实现
type FyneMenuItem struct {
	item    *systray.MenuItem
	handler func()
}

// SetTitle 设置菜单项标题
func (f *FyneMenuItem) SetTitle(title string) {
	if f.item != nil {
		f.item.SetTitle(title)
	}
}

// SetTooltip 设置菜单项提示
func (f *FyneMenuItem) SetTooltip(tooltip string) {
	if f.item != nil {
		f.item.SetTooltip(tooltip)
	}
}

// SetIcon 设置菜单项图标
func (f *FyneMenuItem) SetIcon(iconBytes []byte) {
	if f.item != nil {
		f.item.SetIcon(iconBytes)
	}
}

// Check 勾选菜单项
func (f *FyneMenuItem) Check() {
	if f.item != nil {
		f.item.Check()
	}
}

// Uncheck 取消勾选
func (f *FyneMenuItem) Uncheck() {
	if f.item != nil {
		f.item.Uncheck()
	}
}

// IsChecked 检查是否勾选
func (f *FyneMenuItem) IsChecked() bool {
	if f.item != nil {
		return f.item.Checked()
	}
	return false
}

// Disable 禁用菜单项
func (f *FyneMenuItem) Disable() {
	if f.item != nil {
		f.item.Disable()
	}
}

// Enable 启用菜单项
func (f *FyneMenuItem) Enable() {
	if f.item != nil {
		f.item.Enable()
	}
}

// IsEnabled 检查是否启用
func (f *FyneMenuItem) IsEnabled() bool {
	if f.item != nil {
		return !f.item.Disabled()
	}
	return false
}

// OnClick 设置点击回调
func (f *FyneMenuItem) OnClick(handler func()) {
	if f.item != nil {
		f.handler = handler
		// 重新启动监听goroutine
		go func() {
			for range f.item.ClickedCh {
				if f.handler != nil {
					f.handler()
				}
			}
		}()
	}
}
