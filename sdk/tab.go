package sdk

import (
	"fmt"
	"github.com/package-register/gui/event"
	"image"
	"image/png"
	"os"
	"time"

	"github.com/gonutz/wui/v2"
	"github.com/kbinani/screenshot"
)

// ScreenshotCallback 截图回调函数
type ScreenshotCallback func(img image.Image, err error)

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

// AddImage 添加图片显示组件
func (t *TabContext) AddImage(x, y, w, h int) *ImageDisplay {
	paintBox := wui.NewPaintBox()
	paintBox.SetBounds(x, y, w, h)

	img := &ImageDisplay{
		paintBox: paintBox,
		wuiImage: nil,
		x:        x,
		y:        y,
		width:    w,
		height:   h,
		image:    nil,
	}

	// 设置绘制回调
	paintBox.SetOnPaint(func(canvas *wui.Canvas) {
		if img.image != nil && img.wuiImage != nil {
			// 清空背景
			canvas.FillRect(0, 0, w, h, wui.RGB(240, 240, 240))

			// 获取原始图片尺寸
			srcBounds := img.image.Bounds()
			srcWidth := srcBounds.Dx()
			srcHeight := srcBounds.Dy()

			// 如果图片比显示区域小，直接居中显示
			if srcWidth <= w && srcHeight <= h {
				offsetX := (w - srcWidth) / 2
				offsetY := (h - srcHeight) / 2
				srcRect := wui.Rect(0, 0, srcWidth, srcHeight)
				canvas.DrawImage(img.wuiImage, srcRect, offsetX, offsetY)
				return
			}

			// 计算缩放比例，保持宽高比
			scaleX := float64(w) / float64(srcWidth)
			scaleY := float64(h) / float64(srcHeight)
			scale := scaleX
			if scaleY < scaleX {
				scale = scaleY
			}

			// 计算显示尺寸和居中位置
			displayWidth := int(float64(srcWidth) * scale)
			displayHeight := int(float64(srcHeight) * scale)
			offsetX := (w - displayWidth) / 2
			offsetY := (h - displayHeight) / 2

			// 绘制图片
			srcRect := wui.Rect(0, 0, srcWidth, srcHeight)
			canvas.DrawImage(img.wuiImage, srcRect, offsetX, offsetY)

			// 绘制边框
			canvas.DrawRect(offsetX-1, offsetY-1, displayWidth+2, displayHeight+2, wui.RGB(100, 100, 100))
		} else {
			// 没有图片时显示占位符
			canvas.FillRect(0, 0, w, h, wui.RGB(200, 200, 200))
			canvas.TextRect(0, 0, w, h, "暂无图片", wui.RGB(0, 0, 0))
		}
	})

	// 添加鼠标移动事件处理（用于提示点击）
	paintBox.SetOnMouseMove(func(mouseX, mouseY int) {
		if img.image != nil {
			img.isHover = true
		} else {
			img.isHover = false
		}
	})

	t.panel.Add(paintBox)
	return img
}

// AddScreenshotButton 添加截图按钮
func (t *TabContext) AddScreenshotButton(text string, x, y, w, h int, hideWindow bool, callback ScreenshotCallback) *wui.Button {
	btn := wui.NewButton()
	btn.SetText(text)
	btn.SetBounds(x, y, w, h)

	btn.SetOnClick(func() {
		t.takeScreenshot(hideWindow, callback)
	})

	t.panel.Add(btn)
	return btn
}

// takeScreenshot 执行截图
func (t *TabContext) takeScreenshot(hideWindow bool, callback ScreenshotCallback) {
	// 如果需要隐藏窗口
	if hideWindow && t.app.window != nil {
		originalVisible := t.app.visible
		if originalVisible {
			t.app.HideWindow()
			// 等待窗口完全隐藏
			time.Sleep(200 * time.Millisecond)
		}

		// 延迟截图，确保窗口已隐藏
		go func() {
			time.Sleep(100 * time.Millisecond)
			img, err := t.captureScreen()

			// 恢复窗口显示
			if originalVisible {
				t.app.ShowWindow()
			}

			// 调用回调
			if callback != nil {
				callback(img, err)
			}
		}()
	} else {
		// 直接截图
		go func() {
			img, err := t.captureScreen()
			if callback != nil {
				callback(img, err)
			}
		}()
	}
}

// captureScreen 捕获屏幕
func (t *TabContext) captureScreen() (image.Image, error) {
	// 获取主显示器
	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// ImageDisplay 图片显示组件
type ImageDisplay struct {
	paintBox *wui.PaintBox
	wuiImage *wui.Image
	x        int
	y        int
	width    int
	height   int
	image    image.Image
	mouseX   int
	mouseY   int
	isHover  bool
	onClick  func()
}

// SetImage 设置图片
func (img *ImageDisplay) SetImage(image image.Image) {
	img.image = image
	if image != nil {
		// 转换为wui.Image
		img.wuiImage = wui.NewImage(image)
	} else {
		img.wuiImage = nil
	}
	// 触发重绘
	img.paintBox.Paint()
}

// SetOnClick 设置点击回调
func (img *ImageDisplay) SetOnClick(onClick func()) {
	img.onClick = onClick
}

// GetImage 获取当前图片
func (img *ImageDisplay) GetImage() image.Image {
	return img.image
}

// SaveToFile 保存图片到文件
func (img *ImageDisplay) SaveToFile(filename string) error {
	if img.image == nil {
		return nil
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img.image)
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

// AddChatPanel 添加聊天面板
func (t *TabContext) AddChatPanel(x, y, w, h int) *ChatPanel {
	panel := wui.NewPanel()
	panel.SetBounds(x, y, w, h)
	panel.SetBorderStyle(wui.PanelBorderSunken) // 添加下沉边框，增加深度感

	// 定义内边距
	const padding = 12
	const buttonHeight = 32
	const inputHeight = 36

	// 计算历史显示区域的高度（留出输入区域和按钮空间）
	historyHeight := h - inputHeight - buttonHeight - padding*3

	// 消息历史显示（只读）
	historyEdit := wui.NewTextEdit()
	historyEdit.SetBounds(padding, padding, w-padding*2, historyHeight)
	historyEdit.SetReadOnly(true)
	panel.Add(historyEdit)

	// 输入框区域
	inputY := padding + historyHeight + padding
	inputWidth := w - buttonHeight - padding*3
	inputEdit := wui.NewEditLine()
	inputEdit.SetBounds(padding, inputY, inputWidth, inputHeight)
	panel.Add(inputEdit)

	// 发送按钮
	btnX := padding + inputWidth + padding
	sendBtn := wui.NewButton()
	sendBtn.SetText("发送")
	sendBtn.SetBounds(btnX, inputY, buttonHeight*2, inputHeight) // 稍微宽一点的按钮
	panel.Add(sendBtn)

	chatPanel := &ChatPanel{
		panel:      panel,
		history:    historyEdit,
		input:      inputEdit,
		sendBtn:    sendBtn,
		aiService:  nil,
		onSend:     nil,
		onReceive:  nil,
	}

	// 设置发送按钮点击事件
	sendBtn.SetOnClick(func() {
		if chatPanel.onSend != nil {
			chatPanel.onSend()
		}
	})

	// 注册输入框到应用，用于回车键支持
	t.app.registerChatInput(inputEdit, chatPanel)

	t.panel.Add(panel)
	return chatPanel
}

// ChatPanel 聊天面板组件
type ChatPanel struct {
	panel     *wui.Panel
	history   *wui.TextEdit
	input     *wui.EditLine
	sendBtn   *wui.Button
	aiService *AIService
	onSend    func()
	onReceive  func(message string)
}

// SetAIService 设置 AI 服务
func (c *ChatPanel) SetAIService(aiService *AIService) {
	c.aiService = aiService
}

// OnSend 设置发送回调
func (c *ChatPanel) OnSend(handler func()) {
	c.onSend = handler
}

// OnReceive 设置接收消息回调
func (c *ChatPanel) OnReceive(handler func(message string)) {
	c.onReceive = handler
}

// SendMessage 发送用户消息
func (c *ChatPanel) SendMessage(message string) {
	if message == "" {
		return
	}

	// 显示用户消息
	c.appendMessage("用户", message)

	// 清空输入框
	c.input.SetText("")

	// 如果有 AI 服务，调用 AI
	if c.aiService != nil {
		go func() {
			// 显示"正在生成"提示
			c.appendSystemMessage("AI 正在生成回复...")

			// 添加 AI 消息头
			currentText := c.history.Text()
			timestamp := getCurrentTime()
			header := fmt.Sprintf("\n[%s] AI:\n", timestamp)
			c.history.SetText(currentText + header)
			aiStartPos := len(c.history.Text())

			// 调用 AI 流式接口
			err := c.aiService.ChatStream(message, func(chunk string) {
				// 追加新的内容块
				currentText := c.history.Text()
				c.history.SetText(currentText + chunk)
			})

			if err != nil {
				// 回滚到 AI 消息头之前，添加错误信息
				fullText := c.history.Text()
				c.history.SetText(fullText[:aiStartPos])
				c.appendSystemMessage("❌ AI 调用失败: " + err.Error())
				return
			}

			// 添加换行
			currentText = c.history.Text()
			c.history.SetText(currentText + "\n\n")

			// 触发接收回调
			finalText := c.history.Text()
			if c.onReceive != nil {
				c.onReceive(finalText)
			}
		}()
	}
}

// SendInput 发送当前输入框的内容
func (c *ChatPanel) SendInput() {
	message := c.input.Text()
	c.SendMessage(message)
}

// appendMessage 添加消息到历史记录
func (c *ChatPanel) appendMessage(role, message string) {
	currentText := c.history.Text()
	timestamp := getCurrentTime()
	newMessage := fmt.Sprintf("\n[%s] %s:\n%s\n\n", timestamp, role, message)
	c.history.SetText(currentText + newMessage)

	// 滚动到底部
	// wui.TextEdit 可能没有滚动功能，这里先不处理
}

// appendSystemMessage 添加系统消息
func (c *ChatPanel) appendSystemMessage(message string) {
	currentText := c.history.Text()
	timestamp := getCurrentTime()
	newMessage := fmt.Sprintf("\n[%s] 系统: %s\n\n", timestamp, message)
	_ = currentText // 使用变量避免警告，实际在下一行被使用
	c.history.SetText(currentText + newMessage)
}

// GetHistory 获取聊天历史
func (c *ChatPanel) GetHistory() string {
	return c.history.Text()
}

// ClearHistory 清空聊天历史
func (c *ChatPanel) ClearHistory() {
	c.history.SetText("")
}

// Panel 获取聊天面板（用于添加到 Tab）
func (c *ChatPanel) Panel() *wui.Panel {
	return c.panel
}

// Input 获取输入框（用于高级操作）
func (c *ChatPanel) Input() *wui.EditLine {
	return c.input
}

// History 获取历史文本框（用于高级操作）
func (c *ChatPanel) History() *wui.TextEdit {
	return c.history
}

// getCurrentTime 获取当前时间字符串
func getCurrentTime() string {
	now := time.Now()
	return now.Format("15:04:05")
}
