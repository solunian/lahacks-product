<script>
    import NewPageSvg from "./NewPageSvg.svelte";
    import Links from "./Links.svelte";

    export let text;
    export let message;
    export let count;
    export let oldSearchText;

    const isYou = count % 2 == 0;
    let color = (isYou) ? "bg-teal-50" : "bg-gray-200";
</script>

<div id="stream" class={`${color} h-auto flex-grow p-5 rounded break-words border-y`}>
    {#if isYou}
        <div class="text-right flex flex-col gap-2">
            <h6 class="italic text-gray-500">You</h6>
            {text}
        </div>
    {:else}
        <div class="text-left flex flex-col gap-3">
            <h6 class="italic text-blue-500">Bot</h6>
            <p>
                {text}
            </p>

            <!-- {#each message["links"] as l, i}
                <a href={l["url"]} target="_blank">
                    <button class="text-sm border border-gray-400 py-1 px-3 rounded hover:bg-blue-400 duration-100">
                        {l["name"]} <NewPageSvg/>
                </button></a>
                <p>
                    {l["body"]}
                </p>
            {/each} -->
            {#if (message["links"] && message["links"] != []) } 
                <Links links={message["links"]}/>
            {/if}

            <a href={`https://www.google.com/search?q=${oldSearchText}`} target="_blank">
                <button class="text-md border border-gray-400 py-1 px-3 rounded hover:bg-blue-300 duration-100">
                Check on Google <NewPageSvg/>
            </button></a>
        </div>
    {/if}
    
</div>



<style>
    #stream {
        width: 36rem;
    }
</style>