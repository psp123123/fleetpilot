package api

import (
	"fleetpilot/api/interfaces"
	"fleetpilot/common/logger"
	"fleetpilot/task"
	"sync"
)

// HandlerManager 管理器
type HandlerManager struct {
	handlers map[string]interfaces.ToolHandler
	mutex    sync.RWMutex
}

var (
	manager *HandlerManager
	once    sync.Once
)

// GetHandlerManager 返回单例
func GetHandlerManager() *HandlerManager {
	once.Do(func() {
		manager = &HandlerManager{
			handlers: map[string]interfaces.ToolHandler{
				"nmap": &task.NmapTool{},
				// 新工具在这里手动注册
				// "dirscan": &task.DirScanTool{},
			},
		}
		logger.Debug("已注册工具: %+v", manager.GetAllTools())
	})
	return manager
}

// 根据工具名获取 handler
func (m *HandlerManager) GetHandler(name string) (interfaces.ToolHandler, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	h, ok := m.handlers[name]
	return h, ok
}

// 获取所有工具名称
func (m *HandlerManager) GetAllTools() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	names := make([]string, 0, len(m.handlers))
	for name := range m.handlers {
		names = append(names, name)
	}
	return names
}
