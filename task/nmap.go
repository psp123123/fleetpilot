package task

import (
	"bufio"
	"context"
	"errors"
	"fleetpilot/api"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
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

// 自动注册
func init() {
	api.RegisterTool(&NmapClientParams{})
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

func (n *NmapClientParams) Executed(writer api.WsWriter, msg []byte) error {

	if err := n.PreCheck(); err != nil {
		return err
	}

	args := strings.Fields(n.MsgPayload.ScanParams)
	args = append(args, n.MsgPayload.Target)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "nmap", args...)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return err
	}

	sendLine := func(stream, line string) {
		data := map[string]string{
			"stream": stream,
			"line":   line,
		}
		writer.WriteJSON(data)
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			sendLine("stdout", scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			sendLine("stderr", scanner.Text())
		}
	}()

	if err := cmd.Wait(); err != nil {
		sendLine("stderr", fmt.Sprintf("error: %v", err))
		return err
	}

	writer.WriteJSON(map[string]string{
		"stream": "stdout",
		"line":   "scan done",
	})

	return nil
}
