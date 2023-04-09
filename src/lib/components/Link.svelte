<script>
    import NewPageSvg from "./NewPageSvg.svelte";
    export let label;
    export let url;
</script>


<a href={url}>
    <button class="text-sm border border-gray-400 py-1 px-3 rounded hover:bg-yellow-400 duration-100">
        {label} <NewPageSvg/>
    </button>
</a>