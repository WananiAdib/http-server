package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func write404Response(conn net.Conn) {
	_, err := conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	if err != nil {
		fmt.Println("Error sending 404 request: ", err.Error())
		os.Exit(1)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)

	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		os.Exit(1)
	}

	request := string(buf[:reqLen])
	path := strings.Split(request, " ")
	endpoint := path[1]

	switch {
	case endpoint == "/":
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		if err != nil {
			fmt.Println("Error sending 200 request: ", err.Error())
			os.Exit(1)
		}

	case strings.HasPrefix(endpoint, "/echo"):
		content := endpoint[6:]
		_, err = conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(content), content)))
		if err != nil {
			fmt.Println("Error sending 200 request: ", err.Error())
			os.Exit(1)
		}

	case strings.HasPrefix(endpoint, "/user-agent"):
		lines := strings.Split(request, "\r\n")
		i, err := findUserAgentLine(lines)
		if err != nil {
			fmt.Println("Error reading User-Agent ", err.Error())
			os.Exit(1)
		}
		userAgent := strings.TrimSpace(lines[i][11:])
		_, err = conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent)))
		if err != nil {
			fmt.Println("Error sending 200 request: ", err.Error())
			os.Exit(1)
		}

	case strings.HasPrefix(endpoint, "/files"):
		filepath := os.Args[2] + endpoint[len("/files"):]
		content, err := os.ReadFile(filepath)
		if err != nil {
			write404Response(conn)
			return
		}

		_, err = conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(content), content)))
		if err != nil {
			fmt.Println("Error sending 200 request: ", err.Error())
			os.Exit(1)
		}

	default:
		write404Response(conn)
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}

func findUserAgentLine(lines []string) (int, error) {
	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "User-Agent:") {
			return i, nil
		}
	}

	return -1, fmt.Errorf("User-Agent not found")
}
