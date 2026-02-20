# oAo Agent - Team

一个基于Go语言开发的Windows桌面应用程序，采用SDK风格API设计，支持系统托盘、Tab切换、事件驱动架构和截图功能。

[English](README_EN.md) | 中文

## 特性

- 🎯 **简洁SDK风格API** - 几行代码即可创建完整桌面应用
- 🪟 **系统托盘支持** - 最小化到托盘，后台运行
- 📑 **多Tab界面** - 灵活的页面切换和布局
- 🎨 **中文字体优化** - 默认微软雅黑，支持自定义字体
- ⚡ **事件驱动架构** - 解耦的事件系统，易于扩展
- 🧩 **丰富UI组件** - 按钮、输入框、复选框、进度条等
- 📸 **截图功能** - 支持隐藏窗口截图、图片显示、回调机制
- 🖼️ **图片显示** - 内置图片显示组件，支持保存功能

## 快速开始

### 基本用法

```go
package main

import (
    "gui/sdk"
)

func main() {
    app := sdk.New(
        sdk.WithTitle("我的应用"),
        sdk.WithSize(600, 400),
        sdk.WithTray("我的应用", nil),
        sdk.WithHideConsole(), // 发布时隐藏控制台
    )

    // 注册Tab页面
    app.RegisterTab("主页", func(t *sdk.TabContext) {
        t.AddLabel("欢迎使用", 20, 20, 400, 25)
        t.AddButton("点击我", 20, 60, 100, 30, func() {
            // 按钮点击处理
        })
    })

    // 注册托盘菜单
    app.RegisterTray(func(t *sdk.TrayProxy) {
        t.AddMenuItem("显示/隐藏", "", func() {
            app.ToggleWindow()
        })
        t.AddMenuItem("退出", "", func() {
            app.Exit()
        })
    })

    app.Run()
}
```

### 截图功能

```go
// 添加截图Tab
app.RegisterTab("截图", func(t *sdk.TabContext) {
    // 截图按钮（隐藏窗口）
    t.AddScreenshotButton("截图", 20, 20, 100, 30, true, func(img image.Image, err error) {
        if err != nil {
            log.Printf("截图失败: %v", err)
            return
        }
        log.Printf("截图成功: %dx%d", img.Bounds().Dx(), img.Bounds().Dy())
    })

    // 图片显示区域
    imageDisplay := t.AddImage(20, 60, 400, 300)

    // 截图并显示
    t.AddScreenshotButton("截图并显示", 20, 370, 120, 30, true, func(img image.Image, err error) {
        if err == nil {
            imageDisplay.SetImage(img)
        }
    })
})
```

### 编译运行

```bash
# 开发运行（显示控制台）
go run .

# 发布编译（纯GUI应用）
go build -ldflags "-H=windowsgui" -o myapp.exe
```

## API 参考

### 配置选项

| 选项 | 说明 |
|------|------|
| `WithTitle(title)` | 设置窗口标题 |
| `WithSize(width, height)` | 设置窗口大小 |
| `WithTray(tooltip, icon)` | 启用系统托盘 |
| `WithFont(name, size)` | 设置字体 |
| `WithHideConsole()` | 隐藏控制台（仅编译后有效） |

### UI组件

| 组件 | 方法 | 说明 |
|------|------|------|
| 标签 | `AddLabel(text, x, y, w, h)` | 显示文本 |
| 按钮 | `AddButton(text, x, y, w, h, onClick)` | 可点击按钮 |
| 输入框 | `AddEditLine(x, y, w, h)` | 单行文本输入 |
| 文本框 | `AddTextEdit(x, y, w, h)` | 多行文本输入 |
| 复选框 | `AddCheckBox(text, x, y, w, h, onChange)` | 选择框 |
| 进度条 | `AddProgressBar(x, y, w, h)` | 进度显示 |
| 分隔线 | `AddSeparator(x, y, w)` | 水平分隔线 |
| 图片显示 | `AddImage(x, y, w, h)` | 图片显示组件 |
| 截图按钮 | `AddScreenshotButton(text, x, y, w, h, hideWindow, callback)` | 截图功能 |

### 截图功能

#### ScreenshotCallback
```go
type ScreenshotCallback func(img image.Image, err error)
```

#### ImageDisplay 组件
```go
// 设置图片
imageDisplay.SetImage(img)

// 获取图片
img := imageDisplay.GetImage()

// 保存到文件
err := imageDisplay.SaveToFile("screenshot.png")
```

#### 截图选项
- **hideWindow=false**: 截图时保持窗口显示
- **hideWindow=true**: 截图时自动隐藏窗口，截图后恢复

### 事件系统

```go
app.OnEvent(event.AppStart, func(e event.Event) {
    // 应用启动时
})
app.OnEvent(event.AppExit, func(e event.Event) {
    // 应用退出时
})
app.OnEvent(event.TabSwitch, func(e event.Event) {
    // Tab切换时
})
app.OnEvent(event.WindowShow, func(e event.Event) {
    // 窗口显示时
})
app.OnEvent(event.WindowHide, func(e event.Event) {
    // 窗口隐藏时
})
```

## 技术栈

- **Go 1.21+** - 核心语言
- **wui** - Windows GUI框架
- **fyne.io/systray** - 系统托盘支持
- **kbinani/screenshot** - 截图功能
- **事件驱动架构** - 解耦设计

## 项目结构

```
gui/
├── main.go              # 入口点
├── sdk/
│   ├── gui.go           # 核心SDK
│   ├── tab.go           # Tab上下文 + 截图功能
│   └── tray_proxy.go    # 托盘代理
├── event/
│   └── event.go         # 事件系统
└── tray/
    ├── interface.go     # 托盘接口
    └── fyne_adapter.go  # 托盘适配器
```

## 应用截图

![主界面](screenshots/home.png)
*主界面 - 包含丰富的UI组件*

![截图功能](screenshots/screen.png)
*截图工具 - 支持隐藏窗口截图和图片显示*

![托盘菜单](screenshots/tray.png)
*系统托盘菜单 - 支持窗口控制和快速切换*

![Tab切换](screenshots/tab.png)
*多Tab界面 - 灵活的页面组织*

## 许可证

MIT License
