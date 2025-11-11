package task

// 定义需要前置检查的接口
type PreChecker interface {
	PreCheck() error
}
