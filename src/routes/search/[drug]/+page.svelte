<script>
	export let data;
</script>

<div id="wrapper" class="flex flex-col whitespace-normal gap-4 rounded-lg p-5 border bg-teal-50 mt-16">
	
	<h1 class="text-4xl">{data.name}</h1>

	<div class="flex flex-row gap-2">
		{#each data.med["other_names"].split(", ") as name}
			<div class="border bg-teal-100 rounded-lg p-1 px-2">
				{name}
			</div>
		{/each}
	</div>

	<div>
		<h2 class="text-xl">Uses:</h2>
		{data.med["use"]}
	</div>

	<div>
		<h2 class="text-xl">Side Effects:</h2>
		{data.med["side_effects"]}
	</div>

	<div>
		<h2 class="text-xl">Warnings:</h2>
		{data.med["warnings"]}
	</div>
</div>

<style>
	#wrapper {
		width: 40rem;
	}
</style>
