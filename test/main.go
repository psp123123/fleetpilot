package main

import "fmt"

type NmapInfo struct {
	User  string
	Token string
}

func (n *NmapInfo) PreCheck(u string) string {
	return "precheck NmapInfo done--------%v" + u
}

func (n *NmapInfo) Executed() string {
	return "executed NmapInfo done"
}

func (n *NmapInfo) Success() string {
	return "all NmapInfo success"
}

type Dirsearch struct {
	User  string
	Token string
}

func (n *Dirsearch) PreCheck(u string) string {
	return "precheck Dirsearch done--------%v" + u
}

func (n *Dirsearch) Executed() string {
	return "executed Dirsearch done"
}

func (n *Dirsearch) Success() string {
	return "all Dirsearch success"
}

//ws://192.168.56.11/ws-api/api/ws?tool=nmap&user=admin&token=eyJhbGci

type wsInfo struct {
	Tool  string
	User  string
	Token string
}

// 引入接口
type ToolHandler interface {
	PreCheck(string) string
	Executed() string
	Success() string
}

// 使用接口
func RunTool(t ToolHandler, u string) {
	fmt.Printf("run tool %v,user:%v\n", t.PreCheck(u), u)
	fmt.Printf("run tool %v\n", t.Executed())
	fmt.Printf("run tool %v\n", t.Success())
}

// 注册工具
var toolManger = map[string]ToolHandler{
	"nmap":      &NmapInfo{},
	"dirsearch": &Dirsearch{},
}

func main() {

	ws := wsInfo{
		Tool:  "dirsearch",
		User:  "admin",
		Token: "abc",
	}
	fmt.Printf("get client tool is: %v\n", ws.Tool)

	nm := toolManger[ws.Tool]
	RunTool(nm, ws.User)
}
