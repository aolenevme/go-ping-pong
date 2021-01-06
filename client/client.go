package client

import "syscall/js"

func main() {
	alert := js.Global().Get("alert")
	alert.Invoke("Hello, ping-pong!")
}
