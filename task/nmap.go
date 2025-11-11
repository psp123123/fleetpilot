package task

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
)

// NmapTool 实现 ToolHandler 接口
type NmapTool struct{}

func (n *NmapTool) GetToolName() string {
	return "nmap"
}

// Execute 执行命令并逐行流回前端
func (n *NmapTool) Executed(conn *websocket.Conn, msg []byte) error {
	// 这里可以解析 msg JSON，根据 payload 生成参数
	cmd := exec.CommandContext(context.Background(), "nmap", "-sV", "scanme.nmap.org")

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return err
	}

	sendLine := func(line string) {
		conn.WriteMessage(websocket.TextMessage, []byte(line))
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			sendLine(scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			sendLine(scanner.Text())
		}
	}()

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		if err != nil {
			sendLine(fmt.Sprintf("执行失败: %v", err))
		} else {
			sendLine("任务完成 ✅")
		}
	case <-time.After(30 * time.Second):
		cmd.Process.Kill()
		sendLine("执行超时 ⏱")
	}

	return nil
}
