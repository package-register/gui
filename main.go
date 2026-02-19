package main

import (
	"log"

	"gui/event"
	"gui/sdk"
)

func main() {
	app := sdk.New(
		sdk.WithTitle("oAo Agent - Team"),
		sdk.WithSize(600, 400),
		sdk.WithTray("oAo Agent - Team", nil),
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

	// 注册关于Tab
	app.RegisterTab("关于", func(t *sdk.TabContext) {
		t.AddLabel("oAo Agent - Team", 20, 10, 400, 25)
		t.AddSeparator(20, 40, 540)
		t.AddLabel("版本: 1.0.0", 20, 55, 400, 25)
		t.AddLabel("一个基于Go语言开发的Windows桌面应用程序", 20, 85, 400, 25)
		t.AddLabel("支持系统托盘、Tab切换、事件驱动架构", 20, 115, 400, 25)
		t.AddSeparator(20, 150, 540)
		t.AddLabel("技术栈: Go + wui + fyne.io/systray", 20, 165, 400, 25)
		t.AddLabel("架构: SDK风格 + 事件驱动 + 插拔式托盘", 20, 195, 400, 25)
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
