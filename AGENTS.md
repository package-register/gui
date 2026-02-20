# AGENTS.md - LLM 协作指南

本文档为 LLM（大语言模型）提供 oAo Agent 系统的快速理解和开发指南。

## 📌 项目概览

**oAo Agent** 是一个基于 Go 语言开发的 Windows 桌面应用程序框架，采用 SDK 风格的 API 设计。

**核心特性：**
- SDK 风格 API（Builder 模式 + 函数式选项）
- 事件驱动架构
- 多 Tab 界面
- 系统托盘支持
- 截图功能
- **AI 集成**（基于 trpc-agent-go）

## 📁 关键目录结构

```
gui/
├── main.go              # 应用入口，展示 API 使用
├── sdk/                 # SDK 层（公共 API）
│   ├── gui.go           # App 类，管理应用生命周期
│   ├── tab.go           # TabContext 类，UI 组件和截图
│   ├── ai.go            # AIService 类，AI 对话服务
│   └── tray_proxy.go    # TrayProxy，托盘代理
├── event/               # 事件系统
│   └── event.go         # Event Bus，发布-订阅机制
├── tray/                # 托盘实现
│   ├── interface.go     # Adapter 和 MenuItem 接口
│   └── fyne_adapter.go  # fyne.io/systray 适配器
├── docs/
│   └── SYSTEM_DESIGN.md # ⚠️ 每次修改后必须更新此文件
└── go.mod               # Go 依赖管理
```

## 🎯 设计原则（重要！）

### 1. 分层架构

```
应用层 (main.go)
    ↓
SDK 层 (sdk/)
    ↓
基础设施层 (event/, tray/)
```

**规则：**
- 上层可以依赖下层，下层绝不能依赖上层
- `main.go` 只能调用 `sdk` 包的公共 API
- `sdk` 包可以调用 `event` 和 `tray`
- `event` 和 `tray` 是基础设施，相互独立

### 2. 事件驱动

所有模块通过 `event.Bus` 通信，避免直接依赖。

```go
// 发布事件
app.events.Emit(event.TabSwitch, "Tab名称")

// 订阅事件
app.OnEvent(event.TabSwitch, func(e event.Event) {
    tabName := e.Data.(string)
    log.Printf("切换到: %s", tabName)
})
```

### 3. 适配器模式

托盘功能通过 `Adapter` 接口隔离底层实现，便于替换。

### 4. 延迟初始化

托盘菜单支持在启动前注册，操作被缓存到 `pending` 队列。

## 🔑 核心类说明

### App (sdk/gui.go)

**职责：** 管理应用生命周期、Tab 切换、托盘初始化

**关键方法：**
- `New(opts ...Option) *App` - 创建应用
- `RegisterTab(name, setup)` - 注册 Tab
- `RegisterTray(setup)` - 注册托盘
- `Run() error` - 启动应用（阻塞）
- `SwitchTab(name)` - 切换 Tab
- `ShowWindow()` / `HideWindow()` / `ToggleWindow()` - 窗口控制
- `OnEvent(type, handler)` - 订阅事件

### TabContext (sdk/tab.go)

**职责：** 封装 Tab 页面的 UI 操作

**关键方法：**
- `AddLabel(text, x, y, w, h)` - 添加标签
- `AddButton(text, x, y, w, h, onClick)` - 添加按钮
- `AddEditLine/AddTextEdit` - 添加输入框
- `AddCheckBox(text, x, y, w, h, onChange)` - 添加复选框
- `AddProgressBar(x, y, w, h)` - 添加进度条
- `AddSeparator(x, y, w)` - 添加分隔线
- `AddImage(x, y, w, h)` - 添加图片显示
- `AddScreenshotButton(text, x, y, w, h, hideWindow, callback)` - 添加截图按钮
- `AddChatPanel(x, y, w, h) *ChatPanel` - 添加聊天面板

### Event Bus (event/event.go)

**职责：** 事件订阅和发布

**预定义事件：**
- `AppStart` - 应用启动
- `AppExit` - 应用退出
- `WindowShow` - 窗口显示
- `WindowHide` - 窗口隐藏
- `TabSwitch` - Tab 切换
- `TrayReady` - 托盘就绪

### TrayProxy (sdk/tray_proxy.go)

**职责：** 托盘功能代理

**关键方法：**
- `AddMenuItem(title, tooltip, handler)` - 添加菜单项
- `AddSeparator()` - 添加分隔符
- `SetIcon(icon)` - 设置图标
- `SetTooltip(tooltip)` - 设置提示

### AIService (sdk/ai.go)

**职责：** AI 对话服务封装

**关键方法：**
- `NewAIService(config AIServiceConfig) *AIService` - 创建 AI 服务
- `Chat(message string) (string, error)` - 发送消息并获取回复（同步）
- `ChatStream(message string, callback func(chunk string)) error` - 发送消息并使用流式回调接收回复
- `Close() error` - 关闭 AI 服务

**配置：**
```go
type AIServiceConfig struct {
    APIKey  string  // API Key（必需）
    BaseURL string  // API Base URL（如 OpenAI: https://api.openai.com/v1）
    Model   string  // 模型名称（如 gpt-3.5-turbo）
    UserID  string  // 用户 ID（可选，默认 default-user）
}
```

**依赖：**
- `trpc.group/trpc-go/trpc-agent-go` - LLM Agent 框架
- `trpc.group/trpc-go/trpc-agent-go/model/openai` - OpenAI 兼容客户端
- `trpc.group/trpc-go/trpc-agent-go/runner` - Runner 执行引擎
- `trpc.group/trpc-go/trpc-agent-go/session/inmemory` - 内存会话管理

### ChatPanel (sdk/tab.go)

**职责：** 聊天 UI 组件

**关键方法：**
- `AddChatPanel(x, y, w, h int) *ChatPanel` - 添加聊天面板
- `SetAIService(aiService *AIService)` - 设置 AI 服务
- `SendMessage(message string)` - 发送用户消息（支持流式响应）
- `SendInput()` - 发送当前输入框内容
- `OnSend(handler func())` - 设置发送回调
- `OnReceive(handler func(string))` - 设置接收消息回调
- `GetHistory() string` - 获取聊天历史
- `ClearHistory()` - 清空聊天历史

**特性：**
- 支持流式 AI 响应（实时显示）
- 自动添加时间戳
- 内置错误处理和 UI 反馈
- 异步发送消息，不阻塞 UI

**使用示例：**
```go
chatPanel := t.AddChatPanel(20, 60, 740, 480)
chatPanel.SetAIService(aiService)
chatPanel.OnSend(func() {
    chatPanel.SendInput()
})
```

## 📝 修改代码指南

### 添加新的 UI 组件

**位置：** `sdk/tab.go` 中的 `TabContext`

```go
// 添加新方法
func (t *TabContext) AddDropDown(x, y, w, h int) *wui.DropDown {
    drop := wui.NewDropDown()
    drop.SetBounds(x, y, w, h)
    t.panel.Add(drop)
    return drop
}
```

### 添加新的配置选项

**位置：** `sdk/gui.go`

```go
// 1. 定义 Option 函数
func WithTheme(theme Theme) Option {
    return func(a *App) {
        a.theme = theme
    }
}

// 2. 在 App 结构体添加字段
type App struct {
    // ...
    theme Theme
}
```

### 添加新的事件类型

**位置：** `event/event.go`

```go
const (
    DataChanged Type = "data.changed"  // 新增事件类型
)
```

### 切换托盘底层实现

**位置：** `tray/fyne_adapter.go`

1. 创建新的适配器实现 `CustomAdapter`
2. 实现 `Adapter` 接口的所有方法
3. 在 `tray/interface.go` 中修改 `NewTray()` 使用新适配器

## ⚠️ 重要规则（必须遵守！）

### 1. 依赖规则

✅ 允许：
- `main.go` → `sdk`
- `sdk` → `event`
- `sdk` → `tray`

❌ 禁止：
- `event` → `sdk`
- `tray` → `sdk`
- `event` ↔ `tray`（相互依赖）

### 2. 向后兼容

- 修改公共 API 时保持向后兼容
- 不要删除或重命名已暴露的方法
- 新增功能可以通过新的方法或选项添加

### 3. 错误处理

- 不要使用 `panic` 捕获可预期的错误
- 返回 `error` 让调用者处理
- 在日志中记录关键错误

### 4. 线程安全

- `event.Bus` 已经是线程安全的
- 访问共享状态时使用 `sync.Mutex` 或 `sync.RWMutex`
- UI 操作通常在主线程，注意 goroutine 安全

## 🔍 代码审查检查清单

修改代码时，确认以下几点：

- [ ] 是否遵循了分层架构？
- [ ] 是否使用了事件驱动而非直接依赖？
- [ ] 是否添加了必要的错误处理？
- [ ] 是否更新了 `docs/SYSTEM_DESIGN.md`？（⚠️ 必须做！）
- [ ] 是否有循环依赖？
- [ ] 是否需要更新 `README.md` 示例代码？
- [ ] 是否需要更新 `README_EN.md`？

## 📦 添加新功能的标准流程

1. **理解需求** - 分析功能属于哪个层次
2. **设计接口** - 在合适的包中定义公共接口
3. **实现功能** - 在对应模块中实现
4. **编写测试** - 如有测试框架，添加单元测试
5. **更新文档** - ⚠️ **更新 `docs/SYSTEM_DESIGN.md`**
6. **更新示例** - 如有需要，更新 `main.go` 示例
7. **更新 README** - 如是公共 API，更新文档

## 🔄 每次变更后的必做事项

### 必须更新 `docs/SYSTEM_DESIGN.md`

**更新内容可能包括：**
- 新增的类/方法说明
- 修改的架构图
- 新增的事件类型
- 变更的依赖关系
- 新增的设计模式说明
- 更新最佳实践

**更新的章节：**
- 四、核心模块
- 五、事件驱动机制详解（如有事件变更）
- 六、托盘系统设计（如有托盘变更）
- 七、Tab 系统设计（如有 Tab 变更）
- 十、扩展性设计
- 十一、最佳实践（如有新用法）

### 提交信息规范

使用语义化提交信息：

```
✨ feat: 添加下拉框组件
🐛 fix: 修复截图时窗口隐藏不生效的问题
📝 docs: 更新系统设计文档
⚡ perf: 优化事件发布性能
♻️ refactor: 重构托盘适配器接口
```

## 🚀 快速上手示例

### 创建新 Tab

```go
app.RegisterTab("我的Tab", func(t *sdk.TabContext) {
    t.AddLabel("欢迎", 20, 20, 400, 25)
    t.AddButton("点击我", 20, 60, 100, 30, func() {
        // 按钮点击处理
    })
})
```

### 订阅事件

```go
app.OnEvent(event.TabSwitch, func(e event.Event) {
    tabName := e.Data.(string)
    log.Printf("切换到: %s", tabName)
})
```

### 添加托盘菜单

```go
app.RegisterTray(func(t *sdk.TrayProxy) {
    t.AddMenuItem("显示", "", func() {
        app.ShowWindow()
    })
    t.AddSeparator()
    t.AddMenuItem("退出", "", func() {
        app.Exit()
    })
})
```

## 📚 相关文档

- `docs/SYSTEM_DESIGN.md` - 完整的系统设计文档（⚠️ 每次变更后更新）
- `README.md` - 中文使用文档
- `README_EN.md` - 英文使用文档

## 🤝 与 LLM 协作建议

### 提问时提供上下文

好的提问：
```
我想在 TabContext 中添加一个 DatePicker 组件，应该在哪里添加？需要注意什么？
```

差的提问：
```
怎么添加日期选择器？
```

### 修改任务时

1. 说明要修改哪个文件/模块
2. 说明修改的目的
3. 明确是否需要更新 `docs/SYSTEM_DESIGN.md`
4. 说明是否需要更新 README 示例

### 审查代码时

1. 检查是否违反了分层架构
2. 检查是否缺少事件驱动
3. 检查是否需要更新文档
4. 提供具体的改进建议

## 💡 常见问题

**Q: 可以在 TabContext 中直接访问 App 吗？**

A: 可以，通过 `t.App()` 方法获取 App 实例，但尽量避免这样做，优先使用事件机制通信。

**Q: 如何添加新的底层托盘实现？**

A: 实现 `tray.Adapter` 接口，然后在 `NewTray()` 中使用你的实现。

**Q: 截图功能可以支持指定区域吗？**

A: 当前只支持全屏截图，如需支持区域截图，修改 `captureScreen()` 方法，添加区域参数。

**Q: 如何实现多窗口？**

A: 当前架构是单窗口设计，多窗口需要重构 `App` 类管理多个 `wui.Window` 实例。

---

**最后更新：** 每次代码变更后
**维护者：** LLM + 人工协作
**版本：** 1.0.0
