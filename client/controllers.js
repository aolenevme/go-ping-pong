const DIRECTIONS = Object.freeze({
	RIGHT: "RIGHT",
	LEFT: "LEFT"
});

async function keyDownHandler(e = {}) {
        const key = e.key;
	let direction = DIRECTIONS.RIGHT;

        if (key === "Right" || key === "ArrowRight") {
                direction = DIRECTIONS.RIGHT;
        } else if (key === "Left" || key === "ArrowLeft") {
                direction = DIRECTIONS.LEFT;
        }

        await putData("api/v1/sse", {direction});
}

async function keyUpHandler(e = {}) {
        const key = e.key;

        if (key === "Right" || key === "ArrowRight") {
                rightPressed = false;
        } else if (key === "Left" || key === "ArrowLeft") {
                leftPressed = false;
        }

        await putData("api/v1/sse", {});
}
