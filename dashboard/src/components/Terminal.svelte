<script>
    import { messages } from "../stores/messages";
    import { onDestroy } from "svelte";
    let msgs = [];
    const unsubscribe = messages.subscribe((value) => {
		msgs = value;
	});
	onDestroy(unsubscribe);
</script>

<main>

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
        margin: 0 auto;
    }

    .console {
        background-color: #000;
        margin: 0 auto;
        padding: 1em;
        font-family: "Courier New", monospace;
        height: 800px;
        overflow-y: auto;
    }

</style>
