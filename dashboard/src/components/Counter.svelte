<script>
    import { messages } from "../stores/messages";
    import { onDestroy } from "svelte";
    export let message;
    let count = 0;
    let color = "#fff"
    const unsubscribe = messages.subscribe((value) => {
        count = 0;
        for (let v of value) {
            if (v.Name === message) {
                count = count + 1;
                color = v.Color;
            }
        }
    });
    onDestroy(unsubscribe);
</script>

<main>
    <h3># {message}</h3>
    <div class="wrapper" style="color:{color}">
        {count}
    </div>
</main>

<style>
    main {
        text-align: center;
        padding: 1em;
        margin: 0 auto;
    }

    .wrapper {
        background-color: #000;
        color: #fff;
        margin: 0 auto;
        padding: 1em;
        font-family: "Courier New", monospace;
        font-weight: 500;
        font-size: 60px;
        height: 50px;
        width: 50px;
        
    }
</style>
