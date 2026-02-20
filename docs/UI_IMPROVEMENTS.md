# UI 改进指南

本文档说明如何使用新的 UI 系统，包括布局助手、主题系统和组件样式。

## 目录

1. [快速开始](#快速开始)
2. [布局系统](#布局系统)
3. [主题系统](#主题系统)
4. [改进的组件](#改进的组件)
5. [最佳实践](#最佳实践)
6. [wui 库的限制](#wui-库的限制)
7. [替代方案](#替代方案)

---

## 快速开始

### 基础示例

```go
package main

import (
    "github.com/package-register/gui/sdk"
)

func main() {
    app := sdk.New(
        sdk.WithTitle("我的应用"),
        sdk.WithSize(900, 700),
        sdk.WithTheme(sdk.DefaultTheme()), // 使用默认主题
    )

    app.RegisterTab("主页", setupHomeTab)

    app.Run()
}

func setupHomeTab(t *sdk.TabContext) {
    // 使用主题创建聊天面板
    theme := sdk.DefaultTheme()
    chatPanel := theme.CreateStyledChatPanel(20, 20, 860, 560)

    // 设置 AI 服务
    chatPanel.SetAIService(aiService)

    t.Panel().Add(chatPanel.panel)
}
```

---

## 布局系统

SDK 提供了三种布局类型，类似于前端框架的布局系统：

### 1. 行布局 (Row Layout)

组件水平排列，类似 CSS Flexbox 的 `flex-direction: row`。

```go
func setupRowLayout(t *sdk.TabContext) {
    panel := t.AddPanel(0, 0, 600, 400)

    // 创建行布局
    layout := sdk.NewRowLayout(16, 8, 600, 400, panel)

    // 添加组件（自动水平排列）
    layout.AddButton("按钮1", 100, 40, func() { /* ... */ })
    layout.AddButton("按钮2", 100, 40, func() { /* ... */ })
    layout.AddButton("按钮3", 100, 40, func() { /* ... */ })
}
```

**特点：**
- 组件从左到右排列
- 自动计算 x 坐标
- 使用固定的间距

### 2. 列布局 (Column Layout)

组件垂直排列，类似 CSS Flexbox 的 `flex-direction: column`。

```go
func setupColumnLayout(t *sdk.TabContext) {
    panel := t.AddPanel(0, 0, 400, 600)

    // 创建列布局
    layout := sdk.NewColumnLayout(16, 12, 400, 600, panel)

    // 添加组件（自动垂直排列）
    layout.AddLabel("标题1", 300, 30)
    layout.AddEditLine(300, 35)
    layout.AddLabel("标题2", 300, 30)
    layout.AddEditLine(300, 35)
}
```

**特点：**
- 组件从上到下排列
- 自动计算 y 坐标
- 使用固定的间距

### 3. 网格布局 (Grid Layout)

组件排列成网格，类似 CSS Grid。

```go
func setupGridLayout(t *sdk.TabContext) {
    panel := t.AddPanel(0, 0, 600, 400)

    // 创建 3 列网格布局
    layout := sdk.NewGridLayout(16, 12, 600, 400, 3, panel)

    // 添加组件（自动填入网格）
    for i := 0; i < 6; i++ {
        layout.AddButton(fmt.Sprintf("按钮%d", i+1), 0, 40, nil)
    }
}
```

**特点：**
- 指定列数，自动换行
- 自动计算单元格大小
- 组件均匀分布

### 4. 简易布局函数

对于简单场景，可以使用便捷函数：

```go
// BoxLayout - 垂直或水平排列组件
func exampleBoxLayout(t *sdk.TabContext) {
    panel := t.AddPanel(0, 0, 400, 400)

    controls := []wui.Control{
        wui.NewLabel(),
        wui.NewButton(),
        wui.NewEditLine(),
    }

    widths := []int{300, 100, 200}
    heights := []int{30, 35, 30}

    // 垂直排列
    sdk.BoxLayout(panel, 16, 12, controls, widths, heights, true)

    // 水平排列
    sdk.BoxLayout(panel, 16, 12, controls, widths, heights, false)
}

// GridLayout - 网格排列组件
func exampleGridLayout(t *sdk.TabContext) {
    panel := t.AddPanel(0, 0, 600, 400)

    controls := make([]wui.Control, 6)
    for i := range controls {
        controls[i] = wui.NewButton()
    }

    // 3 列网格
    sdk.GridLayout(panel, 16, 12, 600, 400, 3, controls, nil, nil)
}
```

---

## 主题系统

主题系统提供统一的颜色、字体和间距配置。

### 内置主题

```go
// 浅色主题（Material Design 风格）
theme := sdk.DefaultTheme()

// 深色主题
theme := sdk.DarkTheme()
```

### 自定义主题

```go
customTheme := &sdk.Theme{
    Background:    wui.RGB(240, 240, 245),
    Surface:       wui.RGB(255, 255, 255),
    Primary:       wui.RGB(63, 81, 181),
    Secondary:     wui.RGB(121, 85, 72),
    Accent:        wui.RGB(255, 64, 129),
    DefaultFont:   "微软雅黑",
    FontSize:      -14,
    MediumPadding: 16,
}
```

### 应用主题

```go
app := sdk.New(
    sdk.WithTheme(customTheme),
)
```

### 主题辅助方法

```go
theme := sdk.DefaultTheme()

// 创建带主题样式的组件
panel := theme.CreateStyledPanel(0, 0, 400, 300, wui.PanelBorderSunken)
label := theme.CreateStyledLabel("标题", 16, 16, 300, 30)
btn := theme.CreateStyledButton("点击", 16, 60, 100, 35, func() { /* ... */ })

// 创建样式化的聊天面板
chatPanel := theme.CreateStyledChatPanel(20, 20, 560, 400)

// 获取不同级别的间距
smallPad := theme.GetPadding(1)  // 8px
mediumPad := theme.GetPadding(2) // 16px
largePad := theme.GetPadding(3)  // 24px
```

---

## 改进的组件

### ChatPanel 聊天面板

改进的聊天面板具有更好的布局和间距：

```go
func setupChatPanel(t *sdk.TabContext) {
    // 使用主题创建
    theme := sdk.DefaultTheme()
    chatPanel := theme.CreateStyledChatPanel(20, 20, 740, 480)

    // 设置 AI 服务
    aiService := sdk.NewAIService(sdk.AIServiceConfig{
        APIKey:  "your-api-key",
        BaseURL: "https://api.example.com/v1",
        Model:   "gpt-4",
    })
    chatPanel.SetAIService(aiService)

    // 设置发送回调
    chatPanel.OnSend(func() {
        chatPanel.SendInput()
    })

    // 添加到面板
    t.Panel().Add(chatPanel.panel)
}
```

**改进点：**
- ✅ 支持回车键发送
- ✅ 更好的内边距（12px）
- ✅ 更大的输入框（36px 高度）
- ✅ 更大的按钮（64px 宽）
- ✅ 下沉边框增加深度感

### Tab 标签页

改进的标签页具有更大的可点击区域：

**改进点：**
- ✅ 更大的标签尺寸（120x36px）
- ✅ 更好的间距（8px）
- ✅ 左边距（12px）提升视觉平衡

---

## 最佳实践

### 1. 使用布局系统

✅ **推荐：**
```go
layout := sdk.NewRowLayout(16, 12, 600, 400, panel)
layout.AddButton("按钮1", 100, 40, nil)
layout.AddButton("按钮2", 100, 40, nil)
```

❌ **避免：**
```go
btn1 := wui.NewButton()
btn1.SetBounds(10, 10, 100, 40)
panel.Add(btn1)

btn2 := wui.NewButton()
btn2.SetBounds(120, 10, 100, 40)
panel.Add(btn2)
```

### 2. 使用主题

✅ **推荐：**
```go
theme := sdk.DefaultTheme()
chatPanel := theme.CreateStyledChatPanel(20, 20, 560, 400)
```

❌ **避免：**
```go
chatPanel := t.AddChatPanel(20, 20, 560, 400)
// 手动调整样式...
```

### 3. 使用语义化间距

✅ **推荐：**
```go
x += theme.GetPadding(sdk.SpacingMedium) // 使用主题定义的间距
```

❌ **避免：**
```go
x += 16 // 魔法数字
```

### 4. 代码组织

```go
func setupTab(t *sdk.TabContext) {
    // 1. 创建布局
    layout := sdk.NewColumnLayout(16, 12, 600, 400, t.Panel())

    // 2. 添加组件
    header := layout.AddLabel("标题", 300, 30)
    input := layout.AddEditLine(300, 35)
    btn := layout.AddButton("提交", 120, 40, handleSubmit)

    // 3. 设置事件
    btn.SetOnClick(func() {
        // 处理点击
    })
}
```

---

## wui 库的限制

当前使用的 `github.com/gonutz/wui/v2` 库有以下重要限制：

### 1. 不支持现代布局系统

**限制：**
- 没有内置的 Flexbox 或 Grid 布局
- 只支持绝对定位（x, y, width, height）
- 无法自适应窗口大小

**解决方案：**
- 使用 SDK 提供的布局助手（LayoutHelper）
- 手动计算组件位置

### 2. 样式限制

**限制：**
- 不支持自定义背景色（Panel、Button、Label 等）
- 不支持自定义文字颜色
- 不支持自定义圆角
- 不支持渐变色
- 不支持阴影效果
- 不支持自定义字体颜色

**只能设置：**
- 边框样式（None, SingleLine, Sunken, Raised）
- 字体（名称、大小、粗细）
- 组件位置和大小

### 3. 无法实现 Flutter 级别的视觉效果

由于 wui 是原生 Windows 控件的包装，无法实现：
- 平滑的动画
- Material Design 的涟漪效果
- 自定义绘制
- 半透明效果
- 模糊效果

---

## 替代方案

如果需要更现代的 UI，可以考虑以下替代方案：

### 1. Fyne

**优点：**
- 真正的跨平台（Windows、macOS、Linux）
- 现代化的组件库
- 支持主题系统
- 内置布局系统（HBox、VBox、Grid）
- 良好的文档和社区

**缺点：**
- 相对较新的框架
- 与现有 wui 代码不兼容

**示例：**
```go
import "fyne.io/fyne/v2"

app := app.New()
window := app.NewWindow("Hello")

container := container.NewVBox(
    widget.NewLabel("Hello"),
    widget.NewButton("Click", func() {}),
)

window.SetContent(container)
window.ShowAndRun()
```

### 2. Wails

**优点：**
- 使用 Web 前端技术（HTML/CSS/JavaScript/React/Vue）
- 完全支持现代 UI 效果
- 可以使用任何前端框架
- 丰富的组件库

**缺点：**
- 需要 Web 开发技能
- 应用体积较大
- 与现有 Go 代码集成需要调整

**示例：**
```go
package main

import "github.com/wailsapp/wails/v2"

func main() {
    err := wails.Run(&options.App{
        Title:  "My App",
        Width:  1024,
        Height: 768,
        AssetServer: &assetserver.Options{
            Assets: embed.FS{},
        },
        Background:  context.Background(),
        OnStartup:  app.startup,
        OnDomReady: app.domReady,
        OnShutdown: app.shutdown,
    })
    if err != nil {
        println("Error:", err.Error())
    }
}
```

### 3. Gio

**优点：**
- 纯 Go 实现
- 立即模式渲染
- 高性能
- 支持现代图形效果

**缺点：**
- 学习曲线陡峭
- API 相对底层
- 社区较小

**示例：**
```go
import "gioui.org/app"

func main() {
    go func() {
        w := app.NewWindow()
        var ops op.Ops
        for {
            e := <-w.Events()
            switch e := e.(type) {
            case app.FrameEvent:
                gtx := layout.NewContext(&ops, e)
                // 绘制界面
                e.Frame(gtx.Ops)
            }
        }
    }()
    app.Main()
}
```

### 4. 自定义绘制

可以在 wui 基础上使用 `PaintBox` 组件进行自定义绘制：

```go
paintBox := wui.NewPaintBox()
paintBox.SetBounds(0, 0, 200, 100)
paintBox.SetOnPaint(func(canvas *wui.Canvas) {
    // 绘制自定义内容
    canvas.FillRect(0, 0, 200, 100, wui.RGB(100, 100, 100))
    canvas.TextRect(10, 10, 180, 80, "自定义绘制", wui.RGB(255, 255, 255))
})
```

**注意：** 这需要手动处理所有绘制逻辑，工作量较大。

---

## 总结

### 当前状态

✅ **已实现：**
- 回车键发送消息
- 改进的 ChatPanel 布局
- 更好的 Tab 样式
- 布局助手系统（Row、Column、Grid）
- 主题系统（浅色、深色）
- 语义化间距

### wui 库的限制

⚠️ **无法实现：**
- 自定义组件颜色
- 圆角、阴影等现代视觉效果
- 平滑动画
- Flutter 级别的美观界面

### 建议

1. **短期：** 使用当前的 wui + SDK 改进，适合快速原型和功能开发
2. **中期：** 考虑迁移到 Fyne，获得更好的跨平台和样式支持
3. **长期：** 如果需要最现代的 UI，考虑 Wails（使用 Web 前端）或 Gio（纯 Go 高性能）

---

## 参考资源

- [wui 文档](https://github.com/gonutz/wui)
- [Fyne 文档](https://fyne.io/)
- [Wails 文档](https://wails.io/)
- [Gio 文档](https://gioui.org/)
- [Material Design 规范](https://m3.material.io/)

---

**最后更新：** 2025-02-20
