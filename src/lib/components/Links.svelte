<script>
    import Link from "./Link.svelte";
    export let urls = [];
</script>

<div class="flex flex-col">
    {#each urls as url}
    <Link url={url}/>
    {/each}
</div>