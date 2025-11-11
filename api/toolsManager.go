package api

import (
	"fleetpilot/common/logger"
	"sync"
)

// 定义工具处理客户端传入URL处理过程的接口
type ToolHandler interface {
	// 获取工具名称，用户注册和路由
	GetToolName() string
	Executed(writer WsWriter, msg []byte) error
}

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
		manager = &HandlerManager{handlers: make(map[string]ToolHandler)}
	})
	return manager
}

// 工具注册（由任务包 init() 自动调用）
func RegisterTool(handler ToolHandler) {
	m := GetHandlerManager()
	m.mutex.Lock()
	defer m.mutex.Unlock()

	name := handler.GetToolName()
	m.handlers[name] = handler
	logger.Debug("tool [%s] auto registered", name)
}

// 获取所有已注册的工具名称
func (m *HandlerManager) GetHandler(name string) (ToolHandler, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	h, ok := m.handlers[name]
	return h, ok
}
