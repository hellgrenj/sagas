<script>
	export let message;
	const ws = new WebSocket('ws://localhost:8080/echo');
	ws.onopen = function() {
		ws.send('Hello Server!');
	};
	ws.onmessage = function(e) {
		console.log(e.data);
		message = e.data
	};
</script>

<main>
	<h1>{message}</h1>
</main>

<style>
	main {
		text-align: center;
		padding: 1em;
		max-width: 240px;
		margin: 0 auto;
	}

	h1 {
		color: #ff3e00;
		text-transform: uppercase;
		font-size: 4em;
		font-weight: 100;
	}

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}
</style>