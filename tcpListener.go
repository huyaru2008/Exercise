package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

var form = `<meta charset="utf-8" />                                             
<form method="post">                                                             
用户名:<input name="username" type="text" /><br />                              
密码:<input name="password" type="password" /><br />                            
<input type="submit" /><br />                                                   
</form>                                                                         
 `

func handConn(conn net.Conn) {
	defer conn.Close()
	b, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		return
	}
	fmt.Println(b, err)

	w := bufio.NewWriter(conn)
	defer w.Flush()
	if b.URL.Path == "/" {
		if b.Method == http.MethodGet {
			userCookie, err := b.Cookie("username")
			if err != nil {
				w.WriteString(fmt.Sprintf("HTTP/1.1 200 OK\nContent-Length: %d\nContent-Type: text/html;chatset=utf-8\n\n%s", len(form), form))
			} else {
				str := userCookie.Value
				w.WriteString(fmt.Sprintf("HTTP/1.1 200 OK\nContent-Length: %d\n\n%s", len(str), str))
			}
		} else {
			username := b.FormValue("username")
			// password := b.FormValue("password")
			str := fmt.Sprintf("<h1>%s</h1>", username)
			w.WriteString(fmt.Sprintf("HTTP/1.1 200 OK\nContent-Length: %d\nSet-Cookie: username=%s\n\n%s", len(str), username, str))
		}
	} else {
		w.WriteString("HTTP/1.1 404 Not Found\n\n")
	}

}

func main() {
	tcpListen, err := net.Listen("tcp", ":8080")
	if err != nil {
		return

	}
	fmt.Println("start to listen")
	defer tcpListen.Close()

	for {
		conn, err := tcpListen.Accept()
		if err != nil {
			return
		}
		fmt.Println("accept a connect")
		go handConn(conn)
	}
}
