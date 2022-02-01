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

	const ws = new WebSocket("ws://localhost:8080/ws");
	let prevEv = null;
	let currentColor = getRandomColor();
	ws.onmessage = function (e) {
		const ev = JSON.parse(e.data);
		console.log(ev);
		if (prevEv == null) {
			addNewMsgs({
				Header: true,
				CorrelationId: ev.CorrelationId,
				Color: currentColor,
			});
		}
		if (prevEv && prevEv.CorrelationId !== ev.CorrelationId) {
			currentColor = getRandomColor();
			addNewMsgs({
				Header: true,
				CorrelationId: ev.CorrelationId,
				Color: currentColor,
			});
		}
		ev.Color = currentColor;
		prevEv = ev;
		addNewMsgs(ev);

		msgs = msgs;
	};
	const uniqueCorrelationIds = [];
	function addNewMsgs(ev) {
		if (uniqueCorrelationIds.length == 6) {
			const corrIdToRemove = uniqueCorrelationIds.shift();
			let newMsgs = msgs.filter(
				(m) => m.CorrelationId != corrIdToRemove
			);
			newMsgs.push(ev);
			msgs = newMsgs;
		} else {
			msgs.push(ev);
			msgs = msgs;
		}
		if (!uniqueCorrelationIds.includes(ev.CorrelationId)) {
			uniqueCorrelationIds.push(ev.CorrelationId);
		}
	}
</script>

<main>
	<h3>terminal (latest 5)</h3>
	<div class="console">
		{#each msgs as msg}
			<div style="color: {msg.Color}">
				{#if msg.Header}
					<br />
					{msg.CorrelationId}
				{:else}
					{@html msg.Name}
				{/if}
			</div>
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
		max-width: 600px;
		margin: 0 auto;
		padding: 1em;
		font-family: "Courier New", monospace;
		height: 800px;
		overflow-y: auto;
	}

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}
</style>
