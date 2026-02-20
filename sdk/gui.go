package sdk

import (
	"log"

	"gui/event"
	"gui/tray"

	w32 "github.com/gonutz/w32/v2"
	"github.com/gonutz/wui/v2"
)

// Option 配置选项函数
type Option func(*App)

// TabSetupFunc Tab注册回调
type TabSetupFunc func(t *TabContext)

// TraySetupFunc 托盘注册回调
type TraySetupFunc func(t *TrayProxy)

// App 应用程序主体
type App struct {
	window  *wui.Window
	events  *event.Bus
	tray    *tray.Tray
	visible bool
	font    *wui.Font

	// 配置
	title       string
	width       int
	height      int
	trayEnabled bool
	trayTooltip string
	trayIcon    []byte
	fontName    string
	fontSize    int
	hideConsole bool

	// 注册的回调
	tabSetups map[string]TabSetupFunc
	tabOrder  []string
	traySetup TraySetupFunc
	activeTab string

	// Tab面板
	tabs     map[string]*TabContext
	tabBar   []*wui.Button
	contentY int
}

// New 创建新的GUI应用
func New(opts ...Option) *App {
	app := &App{
		events:    event.NewBus(),
		title:     "oAo Agent",
		width:     600,
		height:    400,
		fontName:  "微软雅黑",
		fontSize:  -14,
		tabSetups: make(map[string]TabSetupFunc),
		tabs:      make(map[string]*TabContext),
		contentY:  40,
	}
	for _, opt := range opts {
		opt(app)
	}
	return app
}

// --- With选项 ---

func WithTitle(title string) Option {
	return func(a *App) { a.title = title }
}

func WithSize(w, h int) Option {
	return func(a *App) { a.width = w; a.height = h }
}

func WithTray(tooltip string, icon []byte) Option {
	return func(a *App) {
		a.trayEnabled = true
		a.trayTooltip = tooltip
		a.trayIcon = icon
	}
}

// WithFont 设置字体
func WithFont(name string, size int) Option {
	return func(a *App) {
		a.fontName = name
		a.fontSize = size
	}
}

// WithHideConsole 隐藏控制台窗口
func WithHideConsole() Option {
	return func(a *App) {
		a.hideConsole = true
	}
}

// --- 注册API ---

// RegisterTab 注册Tab页
func (app *App) RegisterTab(name string, setup TabSetupFunc) {
	app.tabSetups[name] = setup
	app.tabOrder = append(app.tabOrder, name)
}

// RegisterTray 注册托盘菜单
func (app *App) RegisterTray(setup TraySetupFunc) {
	app.traySetup = setup
}

// OnEvent 订阅事件
func (app *App) OnEvent(t event.Type, handler event.Handler) {
	app.events.On(t, handler)
}

// Events 获取事件总线
func (app *App) Events() *event.Bus {
	return app.events
}

// SwitchTab 切换Tab
func (app *App) SwitchTab(name string) {
	if app.activeTab == name {
		return
	}
	// 隐藏当前Tab
	if cur, ok := app.tabs[app.activeTab]; ok {
		cur.hide()
	}
	// 显示新Tab
	if next, ok := app.tabs[name]; ok {
		next.show()
		app.activeTab = name
		app.updateTabBar()
		app.events.Emit(event.TabSwitch, name)
	}
}

// ShowWindow 显示窗口
func (app *App) ShowWindow() {
	if app.window != nil && app.window.Handle() != 0 {
		w32.ShowWindow(w32.HWND(app.window.Handle()), w32.SW_SHOW)
		w32.SetForegroundWindow(w32.HWND(app.window.Handle()))
		app.visible = true
		app.events.Emit(event.WindowShow, nil)
	}
}

// HideWindow 隐藏窗口
func (app *App) HideWindow() {
	if app.window != nil && app.window.Handle() != 0 {
		w32.ShowWindow(w32.HWND(app.window.Handle()), w32.SW_HIDE)
		app.visible = false
		app.events.Emit(event.WindowHide, nil)
	}
}

// ToggleWindow 切换窗口显示/隐藏
func (app *App) ToggleWindow() {
	if app.visible {
		app.HideWindow()
	} else {
		app.ShowWindow()
	}
}

// IsVisible 窗口是否可见
func (app *App) IsVisible() bool {
	return app.visible
}

// Exit 退出应用
func (app *App) Exit() {
	if app.tray != nil {
		app.tray.Quit()
		app.tray = nil // 防止Run()末尾重复Quit
	}
	if app.window != nil {
		app.window.Destroy()
	}
}

// Run 运行应用（阻塞）
func (app *App) Run() error {
	// 创建窗口
	app.window = wui.NewWindow()
	app.window.SetTitle(app.title)
	app.window.SetInnerBounds(100, 50, app.width, app.height)

	// 设置控制台显示
	if app.hideConsole {
		app.window.HideConsoleOnStart()
	}

	// 设置字体
	app.initFont()

	// 构建Tab栏和内容
	app.buildTabBar()
	app.buildTabContents()

	// 窗口关闭行为
	if app.trayEnabled {
		app.window.SetOnCanClose(func() bool {
			app.HideWindow()
			return false
		})
	}

	// 初始化托盘
	if app.trayEnabled {
		app.tray = tray.NewTray()
		if app.trayIcon != nil {
			app.tray.SetIcon(app.trayIcon)
		}
		if app.trayTooltip != "" {
			app.tray.SetTooltip(app.trayTooltip)
		}
		if app.traySetup != nil {
			app.traySetup(&TrayProxy{tray: app.tray})
		}
		if err := app.tray.Start(); err != nil {
			log.Printf("Tray start error: %v", err)
		}
	}

	// 激活默认Tab
	if len(app.tabOrder) > 0 {
		app.activeTab = app.tabOrder[0]
		if t, ok := app.tabs[app.activeTab]; ok {
			t.show()
		}
		app.updateTabBar()
	}

	app.events.Emit(event.AppStart, nil)
	app.visible = true

	// 显示窗口（阻塞直到窗口关闭）
	err := app.window.Show()

	// 窗口关闭后清理托盘
	if app.tray != nil {
		app.tray.Quit()
	}
	app.events.Emit(event.AppExit, nil)

	return err
}

// --- 内部方法 ---

func (app *App) buildTabBar() {
	x := 5
	for _, name := range app.tabOrder {
		tabName := name
		btn := wui.NewButton()
		btn.SetText(tabName)
		btn.SetBounds(x, 5, 80, 28)
		btn.SetOnClick(func() {
			app.SwitchTab(tabName)
		})
		app.window.Add(btn)
		app.tabBar = append(app.tabBar, btn)
		x += 85
	}
}

func (app *App) buildTabContents() {
	for _, name := range app.tabOrder {
		panel := wui.NewPanel()
		panel.SetBounds(0, app.contentY, app.width, app.height-app.contentY)

		ctx := &TabContext{
			name:   name,
			panel:  panel,
			app:    app,
			events: app.events,
		}
		app.tabs[name] = ctx

		if setup, ok := app.tabSetups[name]; ok {
			setup(ctx)
		}

		app.window.Add(panel)
		// 默认隐藏
		panel.SetBounds(0, app.contentY, 0, 0)
	}
}

func (app *App) updateTabBar() {
	for i, name := range app.tabOrder {
		if i < len(app.tabBar) {
			if name == app.activeTab {
				app.tabBar[i].SetText("[ " + name + " ]")
			} else {
				app.tabBar[i].SetText(name)
			}
		}
	}
}

func (app *App) initFont() {
	if app.fontName != "" {
		f, err := wui.NewFont(wui.FontDesc{
			Name:   app.fontName,
			Height: app.fontSize,
		})
		if err != nil && err != wui.NoExactFontMatch {
			log.Printf("Font error: %v, falling back", err)
			return
		}
		app.font = f
		app.window.SetFont(f)
	}
}
