package api

import (
	"fmt"
	"sync"
)

// 注册管理器
type HandlerManager struct {
	handlers map[string]ToolHandler

	// 原子操作
	mutex sync.RWMutex
}

var (
	manager *HandlerManager
	once    sync.Once
)

// 初始化管理器
func GetHandlerManager() *HandlerManager {
	once.Do(func() {
		manager = &HandlerManager{
			handlers: make(map[string]ToolHandler),
		}
	})
	return manager
}

// 注册工具处理器
func (h *HandlerManager) Register(handler ToolHandler) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	name := handler.GetToolName()
	if _, exists := h.handlers[name]; exists {
		panic(fmt.Sprintf("工具 '%s' 的处理器已注册", name))
	}

	h.handlers[name] = handler
}

// 根据工具名称获取处理器
func (h *HandlerManager) GetHandler(toolName string) (ToolHandler, bool) {
	h.mutex.RLock()
	defer h.mutex.Unlock()

	handler, exsits := h.handlers[toolName]
	return handler, exsits
}

// 获取所有已注册的工具名称
func (h *HandlerManager) GetAllHandlers() []string {
	h.mutex.RLock()
	defer h.mutex.Unlock()

	tools := make([]string, 0, len(h.handlers))
	for tool := range h.handlers {
		tools = append(tools, tool)
	}

	return tools
}
