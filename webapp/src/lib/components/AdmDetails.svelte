<script lang="ts">
	let { info }: { info: { properties: Record<string, unknown> } } = $props();
</script>

<div>
	<h1 class="text-2xl font-bold">{info.properties.country}</h1>
	<table class="mt-4">
		<thead>
			<tr class="border-b border-gray-200">
				<th class="pr-2 text-left font-semibold">info</th>
			</tr>
		</thead>
		<tbody>
			{#each Object.entries(info.properties) as [key, value] (key)}
				<tr class="border-b border-gray-200">
					<td class="pr-2">{key.toLowerCase().replaceAll('_', ' ')}</td>
					<td class="p-2">{value}</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>
