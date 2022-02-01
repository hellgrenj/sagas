<script>
	function getRandomColor() {
		var letters = "0123456789ABCDEF";
		var color = "#";
		for (var i = 0; i < 6; i++) {
			color += letters[Math.floor(Math.random() * 16)];
		}
		return color;
	}
	export let msgs = [];
	window.setInterval(() => {
		msgs = []
	}, 10000);
	const ws = new WebSocket("ws://localhost:8080/ws");
	let prevEv = null;
	let currentColor = getRandomColor();
	ws.onmessage = function (e) {
		const ev = JSON.parse(e.data);
		if (prevEv && prevEv.CorrelationId !== ev.CorrelationId) {
			msgs.push({ Name: "<br/>" });
			currentColor = getRandomColor();
		}
		ev.Color = currentColor;
		prevEv = ev;
		msgs.push(ev);

		msgs = msgs;
	};
</script>

<main>
	<h3>terminal (last 10 seconds)</h3>
	<div class="console">
		{#each msgs as msg}
			<div style="color: {msg.Color}">{@html msg.Name}</div>
		{/each}
	</div>
</main>

<style>
	main {
		text-align: center;
		padding: 1em;
		max-width: 240px;
		margin: 0 auto;
	}

	.console {
		background-color: #000;
		max-width: 300px;
		margin: 0 auto;
		padding: 1em;
		font-family: 'Courier New', monospace;
		height: 600px;
		overflow-y: auto;
	}

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}
</style>
