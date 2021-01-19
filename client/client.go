// +build ignore
package main

import (
	"fmt"
	"math"
	"syscall/js"
)

type info struct {
	document      js.Value
	canvas        js.Value
	ctx           js.Value
	interval      js.Value
	canvasWidth   int
	canvasHeight  int
	ballRadius    int
	x             int
	y             int
	dx            int
	dy            int
	paddleWidth   int
	paddleHeight  int
	paddleTopX    int
	paddleBottomX int
	paddleColor   string
	ballColor     string
	rightPressed  bool
	leftPressed   bool
}

func main() {
	document := js.Global().Get("document")
	canvas := document.Call("getElementById", "canvas")
	ctx := canvas.Call("getContext", "2d")
	interval := js.Null()
	canvasWidth := canvas.Get("width").Int()
	canvasHeight := canvas.Get("height").Int()
	ballRadius := 10
	x := canvasWidth / 2
	y := canvasHeight - 30
	dx := 2
	dy := -2
	paddleWidth := 75
	paddleHeight := 10
	paddleTopX, paddleBottomX := (canvasWidth-paddleWidth)/2, (canvasWidth-paddleWidth)/2
	paddleColor := "#141414"
	ballColor := "#d0d0cf"
	rightPressed := false
	leftPressed := false

	i := info{
		document:      document,
		canvas:        canvas,
		ctx:           ctx,
		interval:      interval,
		canvasWidth:   canvasWidth,
		canvasHeight:  canvasHeight,
		ballRadius:    ballRadius,
		x:             x,
		y:             y,
		dx:            dx,
		dy:            dy,
		paddleWidth:   paddleWidth,
		paddleHeight:  paddleHeight,
		paddleTopX:    paddleTopX,
		paddleBottomX: paddleBottomX,
		paddleColor:   paddleColor,
		ballColor:     ballColor,
		rightPressed:  rightPressed,
		leftPressed:   leftPressed,
	}

	drawCb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		draw(&i)

		return nil
	})

	keyDownHandlerCb := createKeyEventListener(&i, keyDownHandler)
	keyUpHandlerCb := createKeyEventListener(&i, keyUpHandler)

	i.interval = js.Global().Call("setInterval", drawCb, 10)
	document.Call("addEventListener", "keydown", keyDownHandlerCb, false)
	document.Call("addEventListener", "keyup", keyUpHandlerCb, false)
	runSse()
	select {}
}

func createKeyEventListener(i *info, keyHandler func(*info, string)) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		key := args[0].Get("key").String()
		keyHandler(i, key)

		return nil
	})
}

func keyDownHandler(i *info, key string) {
	if key == "Right" || key == "ArrowRight" {
		i.rightPressed = true
	} else if key == "Left" || key == "ArrowLeft" {
		i.leftPressed = true
	}
}

func keyUpHandler(i *info, key string) {
	if key == "Right" || key == "ArrowRight" {
		i.rightPressed = false
	} else if key == "Left" || key == "ArrowLeft" {
		i.leftPressed = false
	}
}

func draw(i *info) {
	i.ctx.Call("clearRect", 0, 0, i.canvasWidth, i.canvasHeight)
	drawBall(i)
	drawPaddle(i.paddleTopX, 0, i)
	drawPaddle(i.paddleBottomX, i.canvasHeight-i.paddleHeight, i)

	if i.x > i.canvasWidth-i.ballRadius || i.x < i.ballRadius {
		i.dx = -i.dx
	}

	if shouldRevertBallByY(i) {
		i.dy = -i.dy
	}

	if i.y+i.ballRadius > i.canvasHeight-i.paddleHeight || i.y-i.ballRadius < i.paddleHeight {
		//js.Global().Call("alert", "GAME OVER")
		//i.document.Get("location").Call("reload")
		js.Global().Call("clearInterval", i.interval)
	}

	if i.rightPressed && i.paddleBottomX < i.canvasWidth-i.paddleWidth {
		i.paddleBottomX += 7
	} else if i.leftPressed && i.paddleBottomX > 0 {
		i.paddleBottomX -= 7
	}

	i.x += i.dx
	i.y += i.dy
}

func shouldRevertBallByY(i *info) bool {
	isBallTouchedTopPaddle := i.x >= i.paddleTopX && i.x <= i.paddleTopX+i.paddleWidth && i.y-i.ballRadius <= i.paddleHeight

	isBallTouchedBottomPaddle := i.x >= i.paddleBottomX && i.x <= i.paddleBottomX+i.paddleWidth && i.y+i.ballRadius >= i.canvasHeight-i.paddleHeight

	return isBallTouchedTopPaddle || isBallTouchedBottomPaddle
}

func drawBall(i *info) {
	i.ctx.Call("beginPath")
	i.ctx.Call("arc", i.x, i.y, i.ballRadius, 0, math.Pi*2)
	i.ctx.Set("fillStyle", i.ballColor)
	i.ctx.Call("fill")
	i.ctx.Call("closePath")
}

func drawPaddle(x, y int, i *info) {
	i.ctx.Call("beginPath")
	i.ctx.Call("rect", x, y, i.paddleWidth, i.paddleHeight)
	i.ctx.Set("fillStyle", i.paddleColor)
	i.ctx.Call("fill")
	i.ctx.Call("closePath")
}

func runSse() {
	sse := js.Global().Get("window").Get("EventSource").New("api/v1/sse")

	openCb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("Stream is open")

		return nil
	})

	errorCb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		es := js.Global().Get("window").Get("EventSource")
		rs := args[0].Get("target").Get("readyState").String()

		switch rs {
		case es.Get("CONNECTING").String():
			fmt.Println("Reconnecting...")
		case es.Get("CLOSED").String():
			fmt.Println("Connection failed, will not reconnect")
		default:
			fmt.Println("Unknown error")
		}

		return nil
	})

	msgCb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		data := args[0].Get("data").String()

		JSON := js.Global().Get("window").Get("JSON")

		parsedData := JSON.Call("parse", data)

		fmt.Println(parsedData.Get("Status").String())

		return nil
	})

	sse.Call("addEventListener", "open", openCb)
	sse.Call("addEventListener", "error", errorCb)
	sse.Call("addEventListener", "message", msgCb)
}
