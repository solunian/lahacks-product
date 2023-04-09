<script>
  import Link from "./Link.svelte";
  export let links = [];
  let isExpanded = false;

  function toggleExpansion() {
    isExpanded = !isExpanded;
  }
</script>

<div class="relative">
  <div
    class={`flex flex-col p-5 rounded bg-yellow-50 ${
      isExpanded ? "" : "h-20 overflow-hidden"
    }`}
  >
    <p>More information:</p>
    {#each links as link}
      <Link {link} />
    {/each}
    {#if !isExpanded}
      <div
        class="absolute bottom-0 left-0 w-full h-10 bg-gradient-to-t from-yellow-50 to-transparent"
      />
    {/if}
    {#if isExpanded}
      <button
        class="mt-2 py-1 px-4 bg-yellow-100 hover:bg-yellow-200 rounded-md duration-100 focus:outline-none border-2 border-yellow-400"
        on:click={toggleExpansion}
      >
        Collapse
      </button>
    {/if}
  </div>
  {#if !isExpanded}
    <button
      class="absolute bottom-0 left-1/2 transform -translate-x-1/2 mb-2 py-1 px-4 bg-yellow-100 hover:bg-yellow-200 duration-100 rounded-md focus:outline-none border-2 border-yellow-400"
      on:click={toggleExpansion}
    >
      Expand
    </button>
  {/if}
</div>
