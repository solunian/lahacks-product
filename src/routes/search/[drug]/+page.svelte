<script>
	export let data;
</script>

<div class="flex flex-col whitespace-normal w-3/4 gap-3">
	<h1>{data.name}</h1>

	<div>
		{data.med["other_names"]}
	</div>

	<div>
		{data.med["use"]}
	</div>

	<div>
		{data.med["side_effects"]}
	</div>

	<div>
		{data.med["warnings"]}
	</div>
</div>
