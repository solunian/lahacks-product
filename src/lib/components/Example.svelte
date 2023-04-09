<script>
    export let prompt;
    export let onPick;
</script>


<button class="border w-80 p-4 rounded-lg" on:click={onPick}>
    {`"${prompt}" →`}
</button>