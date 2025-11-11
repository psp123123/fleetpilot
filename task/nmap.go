package task

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
)

// 内部 payload 结构（根据你的实际字段调整）
type NmapPayload struct {
	Host     string `json:"host"`
	ScanType string `json:"scanType"`
	Extra1   string `json:"extra1"`
}

// wrapper 用来解析外层消息
type messageWrapper struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"` // 接受任意 JSON（可能是对象也可能是一个被引号包裹的字符串）
	Extend1 string          `json:"extend1"`
	Extend2 string          `json:"extend2"`
}

// NmapTool 实现 ToolHandler 接口
type NmapTool struct{}

func (n *NmapTool) GetToolName() string {
	return "nmap"
}

// Execute 执行命令并逐行流回前端
func (n *NmapTool) Executed(conn *websocket.Conn, msg []byte) error {
	// 这里可以解析 msg JSON，根据 payload 生成参数
	var w messageWrapper
	if err := json.Unmarshal(msg, &w); err != nil {
		return fmt.Errorf("invalid outer json:%w", err)
	}

	var p NmapPayload

	if len(w.Payload) > 0 && w.Payload[0] == '"' {
		var payloadStr string
		if err := json.Unmarshal(w.Payload, &payloadStr); err != nil {
			return fmt.Errorf("unmarshal payload as string failed: %w", err)
		}
		if err := json.Unmarshal([]byte(payloadStr), &p); err != nil {
			return fmt.Errorf("unmarshal payload string into struct failed: %w", err)
		}
	} else {
		// payload 已经是对象形式，直接解
		if err := json.Unmarshal(w.Payload, &p); err != nil {
			return fmt.Errorf("unmarshal payload object failed: %w", err)
		}
	}

	if p.Host == "" {
		return errors.New("target host empty")
	}

	cmd := exec.CommandContext(context.Background(), "nmap", p.ScanType, p.Host)

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
