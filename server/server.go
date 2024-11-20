package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

var addr = flag.String("addr", "", "TCP服务器监听地址; 默认： \"\" (所有接口).")
var port = flag.Int("port", 8000, "TCP服务器监听端口; 默认: 8000.")

func main() {
	flag.Parse()

	fmt.Println("正在启动服务...")

	src := *addr + ":" + strconv.Itoa(*port)
	listener, _ := net.Listen("tcp", src)
	fmt.Printf("服务监听端口: %s.\n", src)

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("接受客户端连接出错: %s\n", err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("客户端成功连接: " + remoteAddr)

	resp := "连接成功: 服务器当前系统时间: " + time.Now().String() + "\n"
	fmt.Print("< " + resp)
	conn.Write([]byte(resp))

	scanner := bufio.NewScanner(conn)

	for {
		ok := scanner.Scan()

		if !ok {
			break
		}

		handleMessage(scanner.Text(), conn)
	}

	fmt.Println("客户端: " + remoteAddr + " 连接已断开.")
}

func handleMessage(message string, conn net.Conn) {
	fmt.Println("> 请求命令: " + message)

	if len(message) > 0 && message[0] == '/' {
		switch {
		case message == "/time":
			resp := "响应数据: 服务器当前系统时间: " + time.Now().String() + "\n"
			fmt.Print("< " + resp)
			conn.Write([]byte(resp))

		case message == "/quit":
			fmt.Println("正在退出服务.")
			conn.Write([]byte("TCP服务器正在关闭中...\n"))
			fmt.Println("< " + "%关闭成功%")
			conn.Write([]byte("%关闭成功%\n"))
			os.Exit(0)

		default:
			conn.Write([]byte("未注册命令.\n"))
		}
	}
}
