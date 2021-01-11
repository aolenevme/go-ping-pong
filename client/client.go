package main

import "syscall/js"

var (
	math         js.Value = js.Global().Get("Math")
	document     js.Value = js.Global().Get("document")
	canvas       js.Value = document.Call("getElementById", "canvas")
	ctx          js.Value = canvas.Call("getContext", "2d")
	interval     js.Value = js.Null()
	canvasWidth  int      = canvas.Get("width").Int()
	canvasHeight int      = canvas.Get("height").Int()
	ballRadius   int      = 10
	x            int      = canvasWidth / 2
	y            int      = canvasHeight - 30
	dx           int      = 2
	dy           int      = -2
	paddleWidth  int      = 75
	paddleHeight int      = 10
	paddleX      int      = (canvasWidth - paddleWidth) / 2
	paddleColor  string   = "#141414"
	ballColor    string   = "#d0d0cf"
	rightPressed bool     = false
	leftPressed  bool     = false
	isDone                = make(chan bool)
)

func main() {
	interval = js.Global().Call("setInterval", js.FuncOf(draw), 100)
	document.Call("addEventListener", "keydown", js.FuncOf(keyDownHandler), false)
	document.Call("addEventListener", "keyup", js.FuncOf(keyUpHandler), false)
	<-isDone
}

func keyDownHandler(this js.Value, args []js.Value) interface{} {
	event := args[0]
	key := event.Get("key").String()

	if key == "Right" || key == "ArrowRight" {
		rightPressed = true
	} else if key == "Left" || key == "ArrowLeft" {
		leftPressed = true
	}

	return nil
}

func keyUpHandler(this js.Value, args []js.Value) interface{} {
	event := args[0]
	key := event.Get("key").String()

	if key == "Right" || key == "ArrowRight" {
		rightPressed = false
	} else if key == "Left" || key == "ArrowLeft" {
		leftPressed = false
	}

	return nil
}

func draw(this js.Value, args []js.Value) interface{} {
	ctx.Call("clearRect", 0, 0, canvasWidth, canvasHeight)
	drawBall()
	drawPaddle()

	if x > canvasWidth-ballRadius || x < ballRadius {
		dx = -dx
	}

	if (x >= paddleX && x <= paddleX+paddleWidth && y+ballRadius >= canvasHeight-paddleHeight) || y+dy <= ballRadius {
		dy = -dy
	}

	if y+ballRadius > canvasHeight-paddleHeight {
		js.Global().Call("alert", "GAME OVER")
		document.Get("location").Call("reload")
		js.Global().Call("clearInterval", interval)
	}

	if rightPressed && paddleX < canvasWidth-paddleWidth {
		paddleX += 7
	} else if leftPressed && paddleX > 0 {
		paddleX -= 7
	}

	x += dx
	y += dy

	return nil
}

func drawBall() {
	ctx.Call("beginPath")
	ctx.Call("arc", x, y, ballRadius, 0, math.Get("PI").Int()*2)
	ctx.Set("fillStyle", ballColor)
	ctx.Call("fill")
	ctx.Call("closePath")
}

func drawPaddle() {
	ctx.Call("beginPath")
	ctx.Call("rect", paddleX, canvasHeight-paddleHeight, paddleWidth, paddleHeight)
	ctx.Set("fillStyle", paddleColor)
	ctx.Call("fill")
	ctx.Call("closePath")
}
