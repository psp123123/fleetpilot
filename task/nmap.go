package task

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/websocket"
)

// 解析ws消息内容
type NmapClientParams struct {
	MsgType    string     `json:"type"`
	MsgPayload MsgPayload `json:"payload"`
	MsgExtra1  string     `json:"extend1"`
	MsgExtra2  string     `json:"extend2"`
}
type MsgPayload struct {
	Target        string `json:"target"`
	ScanParams    string `json:"scanParams"`
	PayloadExtra1 string `json:"extra1"`
}

// 迁就执行器接口
func (n *NmapClientParams) GetToolName() string {
	return "nmap"
}

// 扫描执行前检测
func (n *NmapClientParams) PreCheck() error {
	// 检查是否为有效的IP地址（IPv4或IPv6）
	if net.ParseIP(n.MsgPayload.Target) == nil {
		return errors.New(" is invalid IP")

		// 检查是否为有效的domain
	} else if !govalidator.IsDNSName(n.MsgPayload.Target) {
		return errors.New(" is invalid domain")
	}

	// 检测扫描参数
	_, hasPrefix := strings.CutPrefix(n.MsgPayload.ScanParams, "-")
	if !hasPrefix && len(n.MsgPayload.ScanParams) > 3 {
		return errors.New("scan params invalid")
	}
	return nil
}

func (n *NmapClientParams) Executed(conn *websocket.Conn, msg []byte) (interface{}, error) {

	// 构建命令参数
	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, n.MsgPayload.ScanParams, n.MsgPayload.Target)

	// 创建带上下文的命令
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "nmap", cmdArgs...)
	output, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("nmap scan timeout")
	}
	return string(output), err
}
