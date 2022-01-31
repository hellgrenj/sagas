<script>
	export let msgs = [];
	const ws = new WebSocket('ws://localhost:8080/ws');
	// ws.onopen = function() {
	// 	ws.send('Hello Server!');
	// };
	ws.onmessage = function(e) {
		console.log(e.data);
		msgs.push(e.data);
		msgs = msgs;
		console.log('length of messages', msgs.length);
	};
</script>

<main>
	<h3>monitoring</h3>
	<div>
	{#each msgs as msg}
		{msg}<br/>
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

	/* .console {
		background-color: #000;
		color:#fff;
		max-width: 300px;
	} */

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}
</style>