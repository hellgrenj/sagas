<script>
	import Terminal from "./components/Terminal.svelte";
	import Chart from "./components/Chart.svelte";
	import { onMount } from "svelte";
	import { connectAndConsume } from "./scripts/ws";
	onMount(connectAndConsume);
</script>

<h1>DASHBOARD</h1>
<div class="top">latest 100 orders</div>
<div class="wrapper">
	<div class="left">
		<Chart />
	</div>
	<div class="right">
		<Terminal />
	</div>
</div>

<style>
	h1 {
		font-family: "Courier New", monospace;
		font-size: 2em;
		margin: 0 auto;
		text-align: center;
		background-color: #000;
		color: #fff;
	}

	.wrapper {
		display: grid;
		grid-template-columns: repeat(3, 1fr);
		gap: 10px;
		grid-auto-rows: minmax(100px, auto);
	}
	.top {
		grid-column: 1/6;
		grid-row: 1;
		text-align: center;
		background-color: #000;
		color: #fff;
	}
	.left {
		grid-column: 1 / 3;
		grid-row: 1;
	}
	.right {
		grid-column: 3 / 6;
		grid-row: 1;
	}
</style>
