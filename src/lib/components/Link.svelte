<script>
  import NewPageSvg from "./NewPageSvg.svelte";
  export let link;

  let domain = new URL(link["url"]).host;

  let isExpanded = false;

  function toggleExpansion() {
    isExpanded = !isExpanded;
  }
</script>

<div class="mb-2">
  <button
    class="flex items-center justify-between w-full py-2 px-4 bg-green-100 hover:bg-green-200 duration-100 rounded-md focus:outline-none"
    on:click={toggleExpansion}
  >
    <span>{link.name} ({domain})</span>
    <svg
      class={`w-5 h-5 transform ${isExpanded ? "rotate-180" : ""}`}
      viewBox="0 0 20 20"
      fill="currentColor"
    >
      <path
        fill-rule="evenodd"
        d="M6.293 7.293a1 1 0 011.414 0L10 9.586l2.293-2.293a1 1 0 011.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z"
        clip-rule="evenodd"
      />
    </svg>
  </button>
  {#if isExpanded}
    <div class="py-2 px-4">
      <p>
        {link["body"]}
      </p>

      <a href={link["url"]} target="_blank">
        <button
          class="text-sm border border-gray-400 py-1 px-3 rounded hover:bg-green-400 duration-100"
        >
          Learn more
          <NewPageSvg />
        </button>
      </a>
    </div>
  {/if}
</div>
