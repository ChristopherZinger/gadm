<script lang="ts">
	import Map from '$lib/components/map/Map.svelte';
	import { mapSelection } from '$lib/stores/map-selection';
</script>

<div style="height: 100vh; width: 100vw; display: flex; flex-direction: row;">
	<div style="flex: 1; padding: 10px;">
		<div
			style="height: 100%; width: 100%; border-radius: 10px; overflow: hidden; border: 1px solid black;"
		>
			<Map />
		</div>
	</div>
	<div class="side-panel">
		{#if $mapSelection}
			<h1 class="text-2xl font-bold">{$mapSelection.properties.country}</h1>
			<table class="mt-4">
				<thead>
					<tr class="border-b border-gray-200">
						<th class="pr-2 text-left font-semibold">info</th>
					</tr>
				</thead>
				<tbody>
					{#each Object.entries($mapSelection.properties) as [key, value] (key)}
						<tr class="border-b border-gray-200">
							<td class="pr-2">{key.toLowerCase().replace('_', ' ')}</td>
							<td class="p-2">{value}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		{/if}
	</div>
</div>

<style>
	.side-panel {
		width: 500px;
		margin: 10px;
		padding: 10px;
	}
</style>
