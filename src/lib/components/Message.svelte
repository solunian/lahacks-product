<script>
    export let text;
    export let count;

    const isYou = count % 2 == 0;
    let color = (isYou) ? "bg-white" : "bg-gray-300";
</script>

<div id="stream" class={`${color} h-auto flex-grow p-5 rounded break-words border-y`}>
    {#if isYou}
        <div class="text-right">
            <h6 class="text-sm italic">You</h6>
            {text}
        </div>
    {:else}
        <div class="text-left">
            <h6 class="text-sm italic">Bot</h6>
            {text}
        </div>
    {/if}
    
</div>



<style>
    #stream {
        width: 40rem;
    }
</style>