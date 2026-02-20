# oAo Agent - Team 系统设计文档

## 一、系统概述

oAo Agent 是一个基于 Go 语言开发的 Windows 桌面应用程序框架，采用 SDK 风格的 API 设计，目标是提供简洁、灵活、易用的桌面应用开发体验。系统支持系统托盘、多 Tab 界面、事件驱动架构和截图功能，适用于快速构建企业级桌面工具应用。

**核心设计目标：**
- 简洁的 API：几行代码即可创建完整的桌面应用
- 解耦架构：各模块职责清晰，易于维护和扩展
- 事件驱动：基于发布-订阅模式，实现松耦合
- 插件化设计：托盘、Tab 等功能模块可独立配置
- 跨平台友好：通过适配器模式支持不同底层实现

---

## 二、设计思想

### 2.1 SDK 风格 API

采用 Builder 模式和函数式选项模式（Functional Options Pattern），让开发者能够通过链式调用和配置函数灵活地创建和配置应用。

```go
app := sdk.New(
    sdk.WithTitle("我的应用"),
    sdk.WithSize(600, 400),
    sdk.WithTray("我的应用", nil),
)
```

**优势：**
- API 简洁直观，学习成本低
- 参数可扩展，不影响现有代码
- 配置与构造分离，代码更清晰

### 2.2 事件驱动架构

整个系统基于事件总线（Event Bus）实现松耦合的组件间通信。核心思想是"发布-订阅"模式，任何组件都可以订阅感兴趣的事件，任何组件都可以发布事件。

**优势：**
- 解耦组件依赖，降低维护成本
- 易于扩展新功能，无需修改现有代码
- 便于实现日志、监控等横切关注点

### 2.3 适配器模式

系统托盘功能采用适配器模式，通过 `Adapter` 接口隔离底层实现细节。当前使用 `fyne.io/systray` 作为底层实现，但可以轻松切换到其他托盘库。

**优势：**
- 底层实现可替换，不影响上层代码
- 便于测试，可以注入 Mock 适配器
- 接口稳定，实现可以演进

### 2.4 延迟初始化策略

托盘菜单支持延迟初始化，即在托盘完全启动之前，可以先调用 API 注册菜单项，这些操作会被缓存到 `pending` 队列中，待托盘启动后统一执行。

**优势：**
- 避免启动时的竞态条件
- API 使用更加灵活，不依赖初始化顺序
- 提升用户体验，托盘启动更快

---

## 三、架构分层

系统采用三层架构设计，从下到上依次为：

```
┌─────────────────────────────────────────────────┐
│                应用层 (main.go)                  │
│  - 应用入口                                      │
│  - Tab 注册                                      │
│  - 托盘菜单注册                                  │
│  - 事件监听                                      │
└─────────────────────────────────────────────────┘
                      ↓
┌─────────────────────────────────────────────────┐
│              SDK 层 (sdk/)                       │
│  - GUI 应用管理 (gui.go)                         │
│  - Tab 上下文 (tab.go)                           │
│  - 托盘代理 (tray_proxy.go)                      │
└─────────────────────────────────────────────────┘
                      ↓
┌─────────────────────────────────────────────────┐
│            基础设施层 (event/, tray/)            │
│  - 事件总线 (event/event.go)                     │
│  - 托盘适配器 (tray/fyne_adapter.go)             │
│  - 托盘接口 (tray/interface.go)                  │
└─────────────────────────────────────────────────┘
```

### 3.1 应用层 (Application Layer)

**职责：**
- 提供应用入口
- 注册 Tab 页面
- 注册托盘菜单
- 订阅应用事件
- 实现业务逻辑

**特点：**
- 只依赖 SDK 层的公共 API
- 不直接操作底层组件
- 通过回调函数实现业务逻辑

### 3.2 SDK 层 (SDK Layer)

**职责：**
- 管理应用生命周期
- 提供 UI 组件的封装
- 管理 Tab 切换
- 托盘功能的代理和封装
- 向应用层暴露简洁的 API

**特点：**
- 封装底层细节，提供高级抽象
- 管理窗口状态和生命周期
- 协调各组件的交互

### 3.3 基础设施层 (Infrastructure Layer)

**职责：**
- 实现事件总线机制
- 提供托盘适配器
- 定义底层接口

**特点：**
- 不依赖上层业务
- 提供通用的基础设施能力
- 可以独立测试和复用

---

## 四、核心模块

### 4.1 应用管理模块 (sdk/gui.go)

**核心类：** `App`

**职责：**
- 管理窗口生命周期（创建、显示、隐藏、销毁）
- 管理 Tab 页面（注册、切换、激活）
- 管理托盘初始化和清理
- 事件总线的持有和分发
- 字体管理

**主要方法：**
- `New(opts ...Option) *App` - 创建应用实例
- `RegisterTab(name string, setup TabSetupFunc)` - 注册 Tab
- `RegisterTray(setup TraySetupFunc)` - 注册托盘
- `Run() error` - 启动应用（阻塞）
- `SwitchTab(name string)` - 切换 Tab
- `ShowWindow()` / `HideWindow()` / `ToggleWindow()` - 窗口控制
- `OnEvent(t event.Type, handler event.Handler)` - 订阅事件

**设计要点：**

1. **配置选项模式**
   使用 `Option` 类型和 `With*` 函数实现灵活的配置：
   ```go
   type Option func(*App)

   func WithTitle(title string) Option {
       return func(a *App) { a.title = title }
   }
   ```

2. **Tab 面板管理**
   每个 Tab 对应一个 `wui.Panel`，通过调整 `Bounds` 来实现显示和隐藏：
   ```go
   func (t *TabContext) show() {
       t.panel.SetBounds(0, t.app.contentY, t.app.width, t.app.height-t.app.contentY)
   }

   func (t *TabContext) hide() {
       t.panel.SetBounds(0, t.app.contentY, 0, 0)
   }
   ```

3. **关闭行为控制**
   启用托盘时，窗口关闭不会退出应用，而是隐藏到托盘：
   ```go
   if app.trayEnabled {
       app.window.SetOnCanClose(func() bool {
           app.HideWindow()
           return false // 阻止窗口关闭
       })
   }
   ```

### 4.2 Tab 上下文模块 (sdk/tab.go)

**核心类：** `TabContext`

**职责：**
- 封装 Tab 页面的 UI 操作
- 提供丰富的 UI 组件添加方法
- 实现截图功能
- 管理图片显示组件

**主要方法：**

**基础组件：**
- `AddLabel(text, x, y, w, h)` - 添加标签
- `AddButton(text, x, y, w, h, onClick)` - 添加按钮
- `AddEditLine(x, y, w, h)` - 添加单行输入框
- `AddTextEdit(x, y, w, h)` - 添加多行文本框
- `AddCheckBox(text, x, y, w, h, onChange)` - 添加复选框
- `AddProgressBar(x, y, w, h)` - 添加进度条
- `AddSeparator(x, y, w)` - 添加分隔线
- `AddPanel(x, y, w, h)` - 添加子面板

**高级功能：**
- `AddImage(x, y, w, h)` - 添加图片显示组件
- `AddScreenshotButton(text, x, y, w, h, hideWindow, callback)` - 添加截图按钮

**设计要点：**

1. **图片自适应显示**
   `ImageDisplay` 组件实现了图片的等比缩放和居中显示：
   - 计算缩放比例：`scale = min(scaleX, scaleY)`
   - 居中偏移：`offset = (容器尺寸 - 显示尺寸) / 2`
   - 绘制边框和占位符

2. **截图功能**
   支持两种截图模式：
   - **隐藏窗口截图**：`hideWindow=true`，截图前隐藏窗口，截图后恢复
   - **普通截图**：`hideWindow=false`，直接截取当前屏幕

   实现细节：
   ```go
   func (t *TabContext) takeScreenshot(hideWindow bool, callback ScreenshotCallback) {
       if hideWindow && t.app.window != nil {
           originalVisible := t.app.visible
           if originalVisible {
               t.app.HideWindow()
               time.Sleep(200 * time.Millisecond) // 等待窗口完全隐藏
           }
           go func() {
               time.Sleep(100 * time.Millisecond)
               img, err := t.captureScreen()
               if originalVisible {
                   t.app.ShowWindow()
               }
               callback(img, err)
           }()
       }
   }
   ```

3. **自定义绘制**
   通过 `wui.PaintBox` 的 `SetOnPaint` 回调实现自定义绘制：
   ```go
   paintBox.SetOnPaint(func(canvas *wui.Canvas) {
       if img.image != nil && img.wuiImage != nil {
           // 绘制图片
           canvas.DrawImage(img.wuiImage, srcRect, offsetX, offsetY)
           // 绘制边框
           canvas.DrawRect(...)
       }
   })
   ```

### 4.3 托盘代理模块 (sdk/tray_proxy.go)

**核心类：** `TrayProxy`

**职责：**
- 作为托盘功能的代理层
- 向应用层暴露简洁的托盘 API
- 隔离底层托盘实现的复杂性

**主要方法：**
- `AddMenuItem(title, tooltip, handler)` - 添加菜单项
- `AddSeparator()` - 添加分隔符
- `SetIcon(icon)` - 设置托盘图标
- `SetTooltip(tooltip)` - 设置提示文本

**设计要点：**
- 代理模式：只暴露必要的接口，隐藏实现细节
- 简化参数：例如省略了 `MenuItem` 的返回值（应用层通常不需要）

### 4.4 事件总线模块 (event/event.go)

**核心类：** `Bus`

**职责：**
- 实现事件的订阅和发布
- 管理事件处理器
- 提供同步和异步事件发布

**主要方法：**
- `On(t Type, h Handler)` - 订阅事件
- `Emit(t Type, data)` - 同步发布事件
- `EmitAsync(t Type, data)` - 异步发布事件

**设计要点：**

1. **线程安全**
   使用 `sync.RWMutex` 保护 `handlers` map：
   ```go
   func (b *Bus) On(t Type, h Handler) {
       b.mu.Lock()
       defer b.mu.Unlock()
       b.handlers[t] = append(b.handlers[t], h)
   }
   ```

2. **事件类型**
   预定义的事件类型：
   ```go
   const (
       AppStart   Type = "app.start"    // 应用启动
       AppExit    Type = "app.exit"     // 应用退出
       WindowShow Type = "window.show"  // 窗口显示
       WindowHide Type = "window.hide"  // 窗口隐藏
       TabSwitch  Type = "tab.switch"   // Tab 切换
       TrayReady  Type = "tray.ready"   // 托盘就绪
   )
   ```

3. **事件传递数据**
   事件可以携带任意类型的数据：
   ```go
   type Event struct {
       EventType Type
       Data      interface{}
   }
   ```

---

## 五、事件驱动机制详解

### 5.1 事件流程

```
事件发布者              事件总线              事件订阅者
    │                      │                      │
    │ ── Emit(event) ──────>│                      │
    │                      │                      │
    │                 遍历 handlers               │
    │                      │                      │
    │                      │ ── handler(event) ──>│
    │                      │                      │
    │                      │ ── handler(event) ──>│
```

### 5.2 典型事件场景

**场景 1：Tab 切换事件**
```go
// 发布：当用户切换 Tab 时
func (app *App) SwitchTab(name string) {
    // 隐藏当前 Tab
    cur.hide()
    // 显示新 Tab
    next.show()
    app.activeTab = name
    app.updateTabBar()
    // 发布事件
    app.events.Emit(event.TabSwitch, name)
}

// 订阅：在 main.go 中
app.OnEvent(event.TabSwitch, func(e event.Event) {
    log.Printf("切换到Tab: %v", e.Data)
})
```

**场景 2：窗口显示/隐藏事件**
```go
// 发布：显示窗口时
func (app *App) ShowWindow() {
    w32.ShowWindow(...)
    app.visible = true
    app.events.Emit(event.WindowShow, nil)
}

// 订阅：在主窗口控制中
app.OnEvent(event.WindowShow, func(e event.Event) {
    // 可以在这里更新 UI 状态
})
```

### 5.3 事件驱动的优势

1. **解耦**：窗口管理、Tab 管理、托盘管理等模块通过事件通信，无需直接引用
2. **可扩展**：新增监听器无需修改现有代码
3. **可测试**：可以注入 Mock 事件总线进行单元测试
4. **可维护**：事件流程清晰，便于理解和调试

---

## 六、托盘系统设计

### 6.1 托盘架构

```
应用层
   ↓
TrayProxy (代理层)
   ↓
Tray (控制层)
   ↓
Adapter (适配器接口)
   ↓
FyneAdapter (实现层)
   ↓
fyne.io/systray (底层库)
```

### 6.2 延迟初始化机制

**问题：** 托盘初始化是异步的，应用启动时可能在托盘就绪前就调用了托盘 API。

**解决方案：** 使用 `pending` 队列缓存操作

```go
type Tray struct {
    adapter Adapter
    running bool
    pending []func() // 待执行的操作队列
}

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

func (t *Tray) Start() error {
    return t.adapter.Initialize(func() {
        t.running = true
        // 执行缓存的操作
        for _, fn := range t.pending {
            fn()
        }
        t.pending = nil
    }, func() {
        t.running = false
    })
}
```

### 6.3 适配器接口设计

```go
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
```

**设计优势：**
- 接口稳定，实现可替换
- 可以轻松切换到其他托盘库
- 便于编写单元测试

### 6.4 FyneAdapter 实现

**关键特性：**
1. **外部循环模式**：使用 `RunWithExternalLoop` 避免阻塞主线程
   ```go
   systray.RunWithExternalLoop(onReady, onExit)
   ```

2. **菜单项监听**：每个菜单项启动一个 goroutine 监听点击事件
   ```go
   go func() {
       for range item.item.ClickedCh {
           handler()
       }
   }()
   ```

3. **线程安全**：使用 `sync.RWMutex` 保护状态

---

## 七、Tab 系统设计

### 7.1 Tab 数据结构

```go
type App struct {
    tabSetups map[string]TabSetupFunc  // Tab 配置函数
    tabOrder  []string                  // Tab 显示顺序
    tabs      map[string]*TabContext   // Tab 实例
    tabBar    []*wui.Button             // Tab 按钮引用
    activeTab string                    // 当前激活的 Tab
}
```

### 7.2 Tab 切换流程

```go
func (app *App) SwitchTab(name string) {
    if app.activeTab == name {
        return
    }
    // 隐藏当前 Tab
    if cur, ok := app.tabs[app.activeTab]; ok {
        cur.hide()  // 设置 Bounds 高度为 0
    }
    // 显示新 Tab
    if next, ok := app.tabs[name]; ok {
        next.show()  // 恢复 Bounds 正常尺寸
        app.activeTab = name
        app.updateTabBar()  // 更新按钮样式
        app.events.Emit(event.TabSwitch, name)  // 发布事件
    }
}
```

### 7.3 Tab 按钮样式更新

```go
func (app *App) updateTabBar() {
    for i, name := range app.tabOrder {
        if i < len(app.tabBar) {
            if name == app.activeTab {
                app.tabBar[i].SetText("[ " + name + " ]")  // 激活状态
            } else {
                app.tabBar[i].SetText(name)  // 普通状态
            }
        }
    }
}
```

### 7.4 Tab 隔离性

每个 Tab 有独立的 `TabContext`，包含：
- 独立的 `wui.Panel`
- 独立的坐标系统
- 独立的 UI 组件引用

**优势：**
- Tab 之间互不干扰
- 便于管理 UI 状态
- 支持 Tab 的独立销毁

---

## 八、截图功能设计

### 8.1 截图架构

```
AddScreenshotButton
    ↓
takeScreenshot(hideWindow, callback)
    ↓
    ├─ hideWindow=true
    │   ├─ HideWindow()
    │   ├─ Sleep(200ms)  // 等待窗口隐藏
    │   ├─ captureScreen()
    │   ├─ ShowWindow()
    │   └─ callback(img, err)
    │
    └─ hideWindow=false
        ├─ captureScreen()
        └─ callback(img, err)
```

### 8.2 屏幕捕获实现

```go
func (t *TabContext) captureScreen() (image.Image, error) {
    // 获取主显示器
    bounds := screenshot.GetDisplayBounds(0)
    // 捕获屏幕
    img, err := screenshot.CaptureRect(bounds)
    return img, err
}
```

### 8.3 图片显示组件

**核心功能：**
1. **等比缩放**：保持宽高比，适应显示区域
2. **居中显示**：图片在容器中居中
3. **占位符**：无图片时显示提示文本
4. **边框绘制**：图片周围绘制灰色边框
5. **保存功能**：支持保存为 PNG 文件
6. **点击交互**：支持设置点击回调

**绘制逻辑：**
```go
paintBox.SetOnPaint(func(canvas *wui.Canvas) {
    if img.image != nil {
        // 计算缩放比例
        scaleX := float64(w) / float64(srcWidth)
        scaleY := float64(h) / float64(srcHeight)
        scale := min(scaleX, scaleY)

        // 计算显示尺寸和居中位置
        displayWidth := int(float64(srcWidth) * scale)
        displayHeight := int(float64(srcHeight) * scale)
        offsetX := (w - displayWidth) / 2
        offsetY := (h - displayHeight) / 2

        // 绘制图片
        canvas.DrawImage(img.wuiImage, srcRect, offsetX, offsetY)

        // 绘制边框
        canvas.DrawRect(offsetX-1, offsetY-1, displayWidth+2, displayHeight+2, wui.RGB(100, 100, 100))
    } else {
        // 显示占位符
        canvas.FillRect(0, 0, w, h, wui.RGB(200, 200, 200))
        canvas.TextRect(0, 0, w, h, "暂无图片", wui.RGB(0, 0, 0))
    }
})
```

---

## 九、依赖关系

### 9.1 核心依赖

| 依赖库 | 版本 | 用途 |
|--------|------|------|
| `github.com/gonutz/wui/v2` | v2.8.2 | Windows GUI 框架，提供窗口、控件等 UI 组件 |
| `github.com/gonutz/w32/v2` | v2.12.1 | Windows API 封装，提供底层 Windows 函数调用 |
| `fyne.io/systray` | v1.12.0 | 跨平台系统托盘库 |
| `github.com/kbinani/screenshot` | latest | 屏幕截图功能库 |

### 9.2 模块依赖图

```
main.go
    ↓
sdk/
    ├─ gui.go ─────┬──→ event/
    │              └──→ tray/
    ├─ tab.go ─────→ event/
    └─ tray_proxy.go ─→ tray/
```

**依赖原则：**
- 单向依赖：上层依赖下层，下层不依赖上层
- 最小依赖：每个模块只依赖必要的其他模块
- 接口隔离：通过接口隔离具体实现

---

## 十、扩展性设计

### 10.1 新增 UI 组件

在 `TabContext` 中添加新方法：

```go
// 添加下拉框
func (t *TabContext) AddComboBox(x, y, w, h int) *wui.ComboBox {
    combo := wui.NewComboBox()
    combo.SetBounds(x, y, w, h)
    t.panel.Add(combo)
    return combo
}
```

### 10.2 新增事件类型

在 `event/event.go` 中定义新类型：

```go
const (
    DataChanged Type = "data.changed"  // 数据变更事件
    NetworkConnected Type = "network.connected"  // 网络连接事件
)
```

### 10.3 切换托盘实现

实现新的适配器：

```go
// 新适配器实现
type CustomAdapter struct {
    // ...
}

func (c *CustomAdapter) Initialize(onReady func(), onExit func()) error {
    // 自定义初始化逻辑
    onReady()
    return nil
}

// 在 NewTray 中替换适配器
func NewTrayWithCustom() *Tray {
    return &Tray{
        adapter: &CustomAdapter{},
    }
}
```

### 10.4 自定义主题

扩展配置选项：

```go
type Theme struct {
    BackgroundColor wui.Color
    TextColor      wui.Color
    AccentColor    wui.Color
}

func WithTheme(theme Theme) Option {
    return func(a *App) {
        a.theme = theme
    }
}
```

---

## 十一、最佳实践

### 11.1 应用启动

```go
func main() {
    app := sdk.New(
        sdk.WithTitle("我的应用"),
        sdk.WithSize(800, 600),
        sdk.WithTray("我的应用", nil),
        sdk.WithHideConsole(),
    )

    // 注册 Tab
    app.RegisterTab("主页", setupHomeTab)
    app.RegisterTab("设置", setupSettingsTab)

    // 注册托盘
    app.RegisterTray(setupTray)

    // 订阅事件
    app.OnEvent(event.AppStart, onAppStart)
    app.OnEvent(event.AppExit, onAppExit)

    // 运行
    if err := app.Run(); err != nil {
        log.Fatalf("应用运行失败: %v", err)
    }
}
```

### 11.2 Tab 页面组织

```go
func setupHomeTab(t *sdk.TabContext) {
    // 添加分隔线
    t.AddSeparator(20, 40, 540)

    // 添加表单控件
    t.AddLabel("用户名:", 20, 55, 60, 25)
    t.AddEditLine(90, 52, 200, 25)

    // 添加按钮
    t.AddButton("提交", 90, 210, 100, 30, onSubmit)
}
```

### 11.3 事件处理

```go
func onAppStart(e event.Event) {
    log.Println("应用启动")
    // 初始化资源
    // 加载配置
}

func onTabSwitch(e event.Event) {
    tabName := e.Data.(string)
    log.Printf("切换到 Tab: %s", tabName)
    // 更新状态
    // 刷新数据
}
```

---

## 十二、总结

oAo Agent 框架通过以下设计原则实现了简洁、灵活、易用的桌面应用开发体验：

1. **分层架构**：应用层、SDK 层、基础设施层清晰分离
2. **事件驱动**：基于发布-订阅模式实现松耦合
3. **适配器模式**：托盘功能底层实现可替换
4. **SDK 风格 API**：简洁直观的 API 设计，学习成本低
5. **延迟初始化**：提升启动性能，避免竞态条件
6. **插件化设计**：Tab、托盘等功能模块可独立配置

该框架为开发 Windows 桌面应用提供了良好的基础，开发者可以专注于业务逻辑，无需关心底层实现细节。
