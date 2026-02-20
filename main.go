package main

import (
	"fmt"
	"image"
	"log"

	"github.com/package-register/gui/event"
	"github.com/package-register/gui/sdk"
)

func main() {
	app := sdk.New(
		sdk.WithTitle("oAo Agent - Team"),
		sdk.WithSize(800, 600),
		sdk.WithTray("oAo Agent - Team", nil),
		sdk.WithHideConsole(), // 隐藏控制台窗口（仅对编译后的exe有效）
	)

	// 注册主页Tab
	app.RegisterTab("主页", func(t *sdk.TabContext) {
		t.AddLabel("欢迎使用 oAo Agent - Team", 20, 10, 400, 25)
		t.AddSeparator(20, 40, 540)

		t.AddLabel("用户名:", 20, 55, 60, 25)
		t.AddEditLine(90, 52, 200, 25)

		t.AddLabel("备注:", 20, 90, 60, 25)
		t.AddTextEdit(90, 87, 200, 80)

		t.AddCheckBox("记住我", 90, 175, 100, 25, func(checked bool) {
			log.Printf("记住我: %v", checked)
		})

		t.AddButton("提交", 90, 210, 100, 30, func() {
			log.Println("提交按钮被点击")
		})

		t.AddLabel("进度:", 20, 260, 60, 25)
		pb := t.AddProgressBar(90, 258, 200, 20)
		pb.SetValue(0.6)
	})

	// 注册截图Tab
	app.RegisterTab("截图", func(t *sdk.TabContext) {
		t.AddLabel("截图工具演示", 20, 10, 400, 25)
		t.AddSeparator(20, 40, 540)

		// 截图按钮区域
		t.AddLabel("截图选项:", 20, 60, 100, 25)

		// 不隐藏窗口截图
		t.AddScreenshotButton("截图（显示窗口）", 20, 90, 150, 30, false, func(img image.Image, err error) {
			if err != nil {
				log.Printf("截图失败: %v", err)
				return
			}
			log.Printf("截图成功，尺寸: %dx%d", img.Bounds().Dx(), img.Bounds().Dy())
		})

		// 隐藏窗口截图
		t.AddScreenshotButton("截图（隐藏窗口）", 180, 90, 150, 30, true, func(img image.Image, err error) {
			if err != nil {
				log.Printf("截图失败: %v", err)
				return
			}
			log.Printf("隐藏窗口截图成功，尺寸: %dx%d", img.Bounds().Dx(), img.Bounds().Dy())
		})

		// 图片显示区域
		t.AddLabel("截图预览:", 20, 140, 100, 25)
		imageDisplay := t.AddImage(20, 170, 400, 250)

		// 图片信息标签
		imageInfoLabel := t.AddLabel("图片信息: 无", 20, 430, 400, 25)

		// 设置点击放大功能
		imageDisplay.SetOnClick(func() {
			if img := imageDisplay.GetImage(); img != nil {
				// 创建放大窗口
				app.RegisterTab("放大查看", func(largeTab *sdk.TabContext) {
					largeTab.AddLabel("图片放大查看", 20, 10, 400, 25)
					largeTab.AddSeparator(20, 40, 760)

					// 显示图片信息
					bounds := img.Bounds()
					info := fmt.Sprintf("原始尺寸: %dx%d", bounds.Dx(), bounds.Dy())
					largeTab.AddLabel(info, 20, 60, 400, 25)

					// 大尺寸图片显示
					largeImage := largeTab.AddImage(20, 90, 760, 420)
					largeImage.SetImage(img)

					// 关闭按钮
					largeTab.AddButton("关闭", 350, 520, 100, 30, func() {
						app.SwitchTab("截图")
					})

					// 保存按钮
					largeTab.AddButton("保存图片", 460, 520, 100, 30, func() {
						if saveErr := largeImage.SaveToFile("screenshot_large.png"); saveErr != nil {
							log.Printf("保存失败: %v", saveErr)
						} else {
							log.Println("图片已保存为: screenshot_large.png")
						}
					})
				})

				// 切换到放大查看Tab
				app.SwitchTab("放大查看")
			}
		})

		// 截图并显示
		t.AddScreenshotButton("截图并显示", 20, 440, 120, 30, true, func(img image.Image, err error) {
			if err != nil {
				log.Printf("截图失败: %v", err)
				return
			}
			log.Printf("截图并显示成功")
			imageDisplay.SetImage(img)

			// 更新图片信息
			bounds := img.Bounds()
			info := fmt.Sprintf("图片尺寸: %dx%d", bounds.Dx(), bounds.Dy())
			imageInfoLabel.SetText(info)
		})

		// 保存截图
		t.AddScreenshotButton("截图并保存", 150, 440, 120, 30, true, func(img image.Image, err error) {
			if err != nil {
				log.Printf("截图失败: %v", err)
				return
			}

			// 生成文件名
			filename := fmt.Sprintf("screenshot_%d.png", img.Bounds().Dx()*img.Bounds().Dy())
			if saveErr := imageDisplay.SaveToFile(filename); saveErr != nil {
				log.Printf("保存失败: %v", saveErr)
			} else {
				log.Printf("截图已保存为: %s", filename)
			}
		})

		// 状态显示
		statusLabel := t.AddLabel("准备就绪", 20, 490, 400, 25)

		// 截图并更新状态
		t.AddScreenshotButton("截图测试", 300, 440, 100, 30, false, func(img image.Image, err error) {
			if err != nil {
				statusLabel.SetText(fmt.Sprintf("截图失败: %v", err))
			} else {
				statusLabel.SetText(fmt.Sprintf("截图成功: %dx%d", img.Bounds().Dx(), img.Bounds().Dy()))
			}
		})
	})

	// 注册关于Tab
	app.RegisterTab("关于", func(t *sdk.TabContext) {
		t.AddLabel("oAo Agent - Team", 20, 10, 400, 25)
		t.AddSeparator(20, 40, 540)
		t.AddLabel("版本: 1.0.0", 20, 55, 400, 25)
		t.AddLabel("一个基于Go语言开发的Windows桌面应用程序", 20, 85, 400, 25)
		t.AddLabel("支持系统托盘、Tab切换、事件驱动架构", 20, 115, 400, 25)
		t.AddSeparator(20, 150, 540)
		t.AddLabel("技术栈: Go + wui + fyne.io/systray + screenshot", 20, 165, 400, 25)
		t.AddLabel("架构: SDK风格 + 事件驱动 + 插拔式托盘", 20, 195, 400, 25)
		t.AddLabel("特性: 截图功能、图片显示、回调机制", 20, 225, 400, 25)
	})

	app.RegisterTab("测试", func(t *sdk.TabContext) {
		t.AddButton("点击", 20, 10, 400, 25, func() {
			log.Println("Hello World!")
		})
	})

	// 注册托盘菜单
	app.RegisterTray(func(t *sdk.TrayProxy) {
		t.AddMenuItem("显示/隐藏", "切换窗口", func() {
			app.ToggleWindow()
		})
		t.AddSeparator()
		t.AddMenuItem("主页", "切换到主页", func() {
			app.SwitchTab("主页")
			app.ShowWindow()
		})
		t.AddMenuItem("截图", "切换到截图", func() {
			app.SwitchTab("截图")
			app.ShowWindow()
		})
		t.AddMenuItem("关于", "切换到关于", func() {
			app.SwitchTab("关于")
			app.ShowWindow()
		})
		t.AddSeparator()
		t.AddMenuItem("退出", "退出程序", func() {
			app.Exit()
		})
	})

	// 事件监听
	app.OnEvent(event.AppStart, func(e event.Event) {
		log.Println("应用已启动")
	})
	app.OnEvent(event.AppExit, func(e event.Event) {
		log.Println("应用已退出")
	})
	app.OnEvent(event.TabSwitch, func(e event.Event) {
		log.Printf("切换到Tab: %v", e.Data)
	})
	app.OnEvent(event.WindowShow, func(e event.Event) {
		log.Println("窗口已显示")
	})
	app.OnEvent(event.WindowHide, func(e event.Event) {
		log.Println("窗口已隐藏")
	})

	// 运行
	if err := app.Run(); err != nil {
		log.Fatalf("应用运行失败: %v", err)
	}
}
