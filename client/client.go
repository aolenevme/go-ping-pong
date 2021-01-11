package main

import "syscall/js"

var (
	math js.Value = js.Global().Get("Math")
	document js.Value = js.Global().Get("document")
	canvas js.Value = document.Call("getElementById", "canvas")
	ctx js.Value = canvas.Call("getContext", "2d")
	interval js.Value = js.Null()
	canvasWidth int = canvas.Get("width").Int()
	canvasHeight int = canvas.Get("height").Int()
	ballRadius int = 10
	x int =  canvasWidth / 2
	y int = canvasHeight - 30
	dx int = 2
	dy int = -2
	paddleWidth int = 75
	paddleHeight int = 10
	paddleX int = (canvasWidth - paddleWidth) / 2
	mainColor string = "#0095DD"
	isDone = make(chan bool)
)

func main() {
	interval = js.Global().Call("setInterval", js.FuncOf(draw), 10)
	<-isDone
}

func draw(this js.Value, args []js.Value) interface{} {
	ctx.Call("clearRect", 0, 0, canvasWidth, canvasHeight)
	drawBall()
	drawPaddle()

	if x + dx > canvasWidth - ballRadius || x + dx < ballRadius {
		dx = -dx
	}

	if y + dy < ballRadius {
		dy = -dy
	} else if y + dy > canvasHeight - ballRadius {
		if x > paddleX && x < paddleX + paddleWidth {
			dy = -dy
		} else {
			js.Global().Call("alert", "GAME OVER")
			document.Get("location").Call("reload")
			js.Global().Call("clearInterval", interval)
		}
	}

	x += dx
	y += dy

	return nil
}

func drawBall() {
	ctx.Call("beginPath")
	ctx.Call("arc", x, y, ballRadius, 0 , math.Get("PI").Int() * 2)
	ctx.Set("fillStyle", mainColor)
	ctx.Call("fill")
	ctx.Call("closePath")
}

func drawPaddle() {
	ctx.Call("beginPath")
	ctx.Call("rect", paddleX, canvasHeight - paddleHeight, paddleWidth, paddleHeight)
	ctx.Set("fillStyle", mainColor)
	ctx.Call("fill")
	ctx.Call("closePath")
}
