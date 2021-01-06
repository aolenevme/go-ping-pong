package main

import "syscall/js"

func main() {
	math := js.Global().Get("Math")
	document := js.Global().Get("document")
	canvas := document.Call("getElementById", "canvas")
	ctx := canvas.Call("getContext", "2d")

	ctx.Call("beginPath")
	ctx.Call("arc", 50, 50, 20, 0, math.Get("PI").Int() * 2)
	ctx.Set("fillStyle", "#0095DD")
	ctx.Call("fill")
	ctx.Call("closePath")
}
