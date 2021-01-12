package main

import "syscall/js"

var (
	math                      = js.Global().Get("Math")
	document                  = js.Global().Get("document")
	canvas                    = document.Call("getElementById", "canvas")
	ctx                       = canvas.Call("getContext", "2d")
	interval                  = js.Null()
	canvasWidth               = canvas.Get("width").Int()
	canvasHeight              = canvas.Get("height").Int()
	ballRadius                = 10
	x                         = canvasWidth / 2
	y                         = canvasHeight - 30
	dx                        = 2
	dy                        = -2
	paddleWidth               = 75
	paddleHeight              = 10
	paddleTopX, paddleBottomX = (canvasWidth - paddleWidth) / 2, (canvasWidth - paddleWidth) / 2
	paddleColor               = "#141414"
	ballColor                 = "#d0d0cf"
	rightPressed              = false
	leftPressed               = false
	isDone                    = make(chan bool)
)

func main() {
	interval = js.Global().Call("setInterval", js.FuncOf(draw), 10)
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
	drawPaddle(paddleTopX, 0)
	drawPaddle(paddleBottomX, canvasHeight-paddleHeight)

	if x > canvasWidth-ballRadius || x < ballRadius {
		dx = -dx
	}

	if shouldRevertBallByY() {
		dy = -dy
	}

	if y+ballRadius > canvasHeight-paddleHeight || y-ballRadius < paddleHeight {
		js.Global().Call("alert", "GAME OVER")
		document.Get("location").Call("reload")
		js.Global().Call("clearInterval", interval)
	}

	if rightPressed && paddleBottomX < canvasWidth-paddleWidth {
		paddleBottomX += 7
	} else if leftPressed && paddleBottomX > 0 {
		paddleBottomX -= 7
	}

	x += dx
	y += dy

	return nil
}

func shouldRevertBallByY() bool {
	isBallTouchedTopPaddle := x >= paddleTopX && x <= paddleTopX+paddleWidth && y-ballRadius <= paddleHeight

	isBallTouchedBottomPaddle := x >= paddleBottomX && x <= paddleBottomX+paddleWidth && y+ballRadius >= canvasHeight-paddleHeight

	return isBallTouchedTopPaddle || isBallTouchedBottomPaddle
}

func drawBall() {
	ctx.Call("beginPath")
	ctx.Call("arc", x, y, ballRadius, 0, math.Get("PI").Int()*2)
	ctx.Set("fillStyle", ballColor)
	ctx.Call("fill")
	ctx.Call("closePath")
}

func drawPaddle(x, y int) {
	ctx.Call("beginPath")
	ctx.Call("rect", x, y, paddleWidth, paddleHeight)
	ctx.Set("fillStyle", paddleColor)
	ctx.Call("fill")
	ctx.Call("closePath")
}
