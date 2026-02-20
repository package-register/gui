package sdk

import (
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
