<script lang="ts">
	import AdmDetails from '$lib/components/AdmDetails.svelte';
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
		<div>
			<div>
				<h1 class="text-2xl font-bold">Hello World Lines!</h1>
				<p>
					Explore worldwide administrative boundaries on an interactive map. Click anywhere on the
					map to pick a unit and see its details here in the side panel.
				</p>
			</div>
			{#if $mapSelection && $mapSelection.type === 'adm'}
				<div class="py-4">
					<AdmDetails info={$mapSelection} />
				</div>
			{/if}
		</div>
		<div>
			<a class="underline" href="https://docs.worldlines.dev" target="_blank">API docs</a>
		</div>
	</div>
</div>

<style>
	.side-panel {
		width: 500px;
		margin: 10px;
		padding: 10px;
		display: flex;
		flex-direction: column;
		gap: 10px;
		justify-content: space-between;
	}
</style>
