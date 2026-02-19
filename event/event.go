package event

import "sync"

// Type 事件类型
type Type string

const (
	AppStart   Type = "app.start"
	AppExit    Type = "app.exit"
	WindowShow Type = "window.show"
	WindowHide Type = "window.hide"
	TabSwitch  Type = "tab.switch"
	TrayReady  Type = "tray.ready"
)

// Event 事件
type Event struct {
	EventType Type
	Data      interface{}
}

// Handler 事件处理函数
type Handler func(Event)

// Bus 事件总线
type Bus struct {
	mu       sync.RWMutex
	handlers map[Type][]Handler
}

// NewBus 创建事件总线
func NewBus() *Bus {
	return &Bus{handlers: make(map[Type][]Handler)}
}

// On 订阅事件
func (b *Bus) On(t Type, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[t] = append(b.handlers[t], h)
}

// Emit 发布事件
func (b *Bus) Emit(t Type, data interface{}) {
	b.mu.RLock()
	handlers := b.handlers[t]
	b.mu.RUnlock()
	e := Event{EventType: t, Data: data}
	for _, h := range handlers {
		h(e)
	}
}

// EmitAsync 异步发布事件
func (b *Bus) EmitAsync(t Type, data interface{}) {
	go b.Emit(t, data)
}
