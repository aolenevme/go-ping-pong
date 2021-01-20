const canvas = document.getElementById("canvas");
const ctx = canvas.getContext("2d");
const ballRadius = 10;
const paddleWidth = 75;
const paddleHeight = 10;
const paddleColor = "#141414";
const ballColor = "#d0d0cf";


let x = canvas.width / 2;
let y = canvas.height - 30;
let dx = 2;
let dy = -2;
let paddleTopX = (canvas.width - paddleWidth) / 2;
let paddleBottomX = (canvas.width - paddleWidth) / 2;
let rightPressed = false;
let leftPressed = false;

let interval = setInterval(draw, 10);
document.addEventListener("keydown", keyDownHandler, false);
document.addEventListener("keyup", keyUpHandler, false);
runSse();

function draw() {
	ctx.clearRect(0, 0, canvas.width, canvas.height);
	drawBall();
	drawPaddle(paddleTopX, 0);
	drawPaddle(paddleBottomX, canvas.height - paddleHeight)

	if ((x > canvas.width - ballRadius) || x < ballRadius) {
		dx = -dx;
	}

	if (shouldRevertBallByY()) {
		dy = -dy;
	}

	if ((y + ballRadius > canvas.height - paddleHeight) || (y - ballRadius) < paddleHeight) {
		alert("Game Over");
		document.location.reload();
		clearInterval(interval);
	}

	if (rightPressed && (paddleBottomX < canvas.width - paddleWidth)) {
		paddleBottomX += 7;
	} else if (leftPressed && paddleBottomX > 0) {
		paddleBottomX -= 7;
	}

	x += dx;
	y += dy;
}

function shouldRevertBallByY() {
	const isBallTouchedTopPaddle = x >= paddleTopX && x < paddleTopX + paddleWidth && y - ballRadius <= paddleHeight;
	const isBallTouchedBottomPaddle = x >= paddleBottomX && x < paddleBottomX + paddleWidth && y + ballRadius >= canvas.height - paddleHeight;

	return isBallTouchedTopPaddle || isBallTouchedBottomPaddle;
}

function drawBall() {
	ctx.beginPath();
	ctx.arc(x, y, ballRadius, 0, Math.PI * 2);
	ctx.fillStyle = ballColor;
	ctx.fill();
	ctx.closePath();
}

function drawPaddle(x = 0, y = 0) {
	ctx.beginPath();
	ctx.rect(x, y, paddleWidth, paddleHeight);
	ctx.fillStyle = paddleColor;
	ctx.fill();
	ctx.closePath();
}

function keyDownHandler(e = {}) {
	const key = e.key;

	if (key === "Right" || key === "ArrowRight") {
		rightPressed = true;
	} else if (key === "Left" || key === "ArrowLeft") {
		leftPressed = true;
	}
}

function keyUpHandler(e = {}) {
	const key = e.key;

	if (key === "Right" || key === "ArrowRight") {
		rightPressed = false;
	} else if (key === "Left" || key === "ArrowLeft") {
		leftPressed = false;
	}
}

function runSse() {
	const sse = new EventSource("api/v1/sse");
	
	sse.addEventListener("open", () => {
		console.log("Stream is open");
	});

	sse.addEventListener("error", e => {
		const eventSourceState = e.target.readyState;

		switch (eventSourceState) {
			case EventSource.CONNECTING: 
				console.log("Reconnecting");
				break;
			case EventSource.CLOSED:
				console.log("Connections failed, will not reconnect");
				break;
			default:
				console.log("Unknown error");
		}
	});

	sse.addEventListener("message", e => {
		const data = JSON.parse(e.data);
		
		console.dir(data);
	});
}
