<script>
    import { messages } from "../stores/messages";
    import { onDestroy } from "svelte";
    let msgs = [];
    const unsubscribe = messages.subscribe((value) => {
        console.log(value)
		msgs = value;
	});
	onDestroy(unsubscribe);
</script>

<main>
    <h3>Latest 5 orders</h3>
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
