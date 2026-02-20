package sdk

import "github.com/package-register/gui/tray"

// TrayProxy 托盘代理，暴露给用户的简洁API
type TrayProxy struct {
	tray *tray.Tray
}

// AddMenuItem 添加菜单项
func (p *TrayProxy) AddMenuItem(title, tooltip string, handler func()) {
	p.tray.AddMenuItem(title, tooltip, handler)
}

// AddSeparator 添加分隔符
func (p *TrayProxy) AddSeparator() {
	p.tray.AddSeparator()
}

// SetIcon 设置图标
func (p *TrayProxy) SetIcon(icon []byte) {
	p.tray.SetIcon(icon)
}

// SetTooltip 设置提示
func (p *TrayProxy) SetTooltip(tooltip string) {
	p.tray.SetTooltip(tooltip)
}
