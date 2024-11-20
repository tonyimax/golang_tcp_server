package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"
)

var host = flag.String("host", "localhost", "The hostname or IP to connect to; defaults to \"localhost\".")
var port = flag.Int("port", 8000, "The port to connect to; defaults to 8000.")
var isNeedReConnect = flag.Bool("need-reconnect", false, "Whether to need a reconnect")
var gConn net.Conn
var reader *bufio.Reader

func main() {
	flag.Parse()

	dest := *host + ":" + strconv.Itoa(*port)
	fmt.Printf("正在连接到TCP服务器: %s...\n", dest)

	conn, err := net.Dial("tcp", dest)

	if err != nil {
		if _, t := err.(*net.OpError); t {
			fmt.Println("连接到TCP服务器出错.")
		} else {
			fmt.Println("未知错误: " + err.Error())
		}
		fmt.Println("===休眠1秒，重连TCP服务器===", time.Now().String())
		time.Sleep(1 * time.Second)
		gConn = reConection()
		if nil != gConn {
			*isNeedReConnect = false
			go readConnection(gConn)
		}
	}

	gConn = conn

	go readConnection(gConn)

	for {
		fmt.Println("===循环处理数据===")
		reader = bufio.NewReader(os.Stdin)
		if *isNeedReConnect {
			gConn = reConection()
			if nil != gConn {
				*isNeedReConnect = false
				fmt.Println("重新连接TCP服务器后，启动新协程处理连接数据===")
				go readConnection(gConn)
			} else {
				fmt.Println("===休眠1秒，重连TCP服务器===", time.Now().String())
				time.Sleep(1 * time.Second)
				gConn = reConection()
				if nil != gConn {
					*isNeedReConnect = false
					go readConnection(gConn)
				}
			}
		}
		if nil != gConn {
			fmt.Print("> ")
			text, _ := reader.ReadString('\n')
			fmt.Println("===>gConn: ", gConn)
			gConn.SetWriteDeadline(time.Now().Add(1 * time.Second)) //设置写入超时1秒
			_, err := gConn.Write([]byte(text))
			if err != nil {
				fmt.Println("写入网络流失败.")
				fmt.Print("===gConn is nil ===")
				fmt.Println("===休眠1秒，重连TCP服务器===", time.Now().String())
				time.Sleep(1 * time.Second)
				gConn = reConection()
				if nil != gConn {
					*isNeedReConnect = false
					go readConnection(gConn)
				}
				//break
			}
		} else {
			fmt.Print("===gConn is nil ===")
			fmt.Println("===休眠1秒，重连TCP服务器===", time.Now().String())
			time.Sleep(1 * time.Second)
			gConn = reConection()
			if nil != gConn {
				*isNeedReConnect = false
				go readConnection(gConn)
			}
		}

	}

}

func reConection() net.Conn {
	dest := *host + ":" + strconv.Itoa(*port)
	fmt.Printf("正在重新连接到TCP服务器: %s...\n", dest)

	conn, err := net.Dial("tcp", dest)

	if err != nil {
		if _, t := err.(*net.OpError); t {
			fmt.Println("重新连接到TCP服务器出错.")
		} else {
			fmt.Println("未知错误: " + err.Error())
		}
		return nil
	}
	*isNeedReConnect = false
	return conn
}

func readConnection(conn net.Conn) {
	t1 := 0
	for {
		if nil == conn {
			break
		}
		scanner := bufio.NewScanner(conn)
		isQuit := false
		for {
			ok := scanner.Scan()
			text := scanner.Text()

			command := handleCommands(text)
			if !command {
				fmt.Printf("\b\b** %s\n> ", text)
			}

			if !ok {
				fmt.Println("读取TCP服务器连接失败")
				if t1 == 0 {
					t1 = time.Now().Second()
				} else {
					fmt.Printf("===> t1: %v\n", t1)
					fmt.Printf("===> Now.Second: %v\n", time.Now().Second())
					if time.Now().Second()-t1 >= 5 {
						isQuit = true
						fmt.Println("连续5秒连接不上服务器，需要退出协程并重新连接服务器")
						*isNeedReConnect = true
						t1 = 0
					}
				}
				break
			}
		}
		if isQuit {
			fmt.Println("连续5秒连接不上服务器，正在退出协程（回车后重新连接服务）===")
			fmt.Print("> ")
			main()
			break
		}
	}

}

func handleCommands(text string) bool {
	r, err := regexp.Compile("^%.*%$")
	if err == nil {
		if r.MatchString(text) {

			switch {
			case text == "%关闭成功%":
				fmt.Println("\b\bTCP服务已关闭，将断开连接")
				os.Exit(0)
			}

			return true
		}
	}

	return false
}
