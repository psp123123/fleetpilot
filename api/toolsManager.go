package api

import (
	"fleetpilot/common/logger"
	"fleetpilot/task"
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
		toolMap := make(map[string]ToolHandler)
		toolMap["nmap"] = &task.NmapClientParams{}
		manager = &HandlerManager{
			handlers: toolMap,
		}
	})
	logger.Debug("注册的信息：%v", manager)
	return manager
}

// 根据工具名称获取处理器
func (h *HandlerManager) GetHandler(toolName string) (ToolHandler, bool) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	handler, exsits := h.handlers[toolName]
	logger.Debug("根据工具名称获取处理器-%v", handler)
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
