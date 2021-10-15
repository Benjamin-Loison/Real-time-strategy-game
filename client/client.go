package main
import (
	"fmt"
	"net"
	"os"
	"bufio"
	"strings"
)

func write(conn *net.Conn, keep_going *bool, ch chan string) {
	reader := bufio.NewReader(os.Stdin)
	for *keep_going {
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')

		if (strings.Compare(text, "/exit\n") == 0) {
			*keep_going = false
			close(ch)
		} else {
			fmt.Fprintf(*conn, text)
		}
	}
}
func read (conn *net.Conn, keep_going *bool, ch chan string) {
	p :=  make([]byte, 2048)
	for *keep_going {
		_, err := bufio.NewReader(*conn).Read(p)
		if err == nil {
			fmt.Printf("%s\n", p)
		} else {
			fmt.Printf("Some error %v\n", err)
		}
	}
	close(ch)
}

func main() {
	conn, err := net.Dial("udp", "192.168.9.128:8000")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}

	keep_going := true

	in_stdin := make(chan string)
	in_server := make(chan string)

	go read(&conn, &keep_going, in_stdin)
	go write(&conn, &keep_going, in_server)

	for keep_going {
		select {
		case s1 := <-in_stdin :
			fmt.Println (s1)
		case s2 := <-in_server:
			fmt.Println(s2)
		}
	}

	conn.Close()
}

