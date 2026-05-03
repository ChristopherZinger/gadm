<script lang="ts">
	import AdmDetails from '$lib/components/AdmDetails.svelte';
	import CodeEditorInput from '$lib/components/CodeEditorInput.svelte';
	import Map from '$lib/components/map/Map.svelte';
	import type { SidePanelView } from '$lib/components/nav/Nav.svelte';
	import Nav from '$lib/components/nav/Nav.svelte';
	import { mapSelection } from '$lib/stores/map-selection';

	let sidePanelView = $state<SidePanelView>('geojson');
</script>

<div class="flex h-screen w-screen flex-row gap-3 p-3">
	<div style="flex-grow: 1; flex-shrink: 1;">
		<div class="map h-full w-full overflow-hidden rounded-md">
			<Map />
		</div>
	</div>
	<div class="side-panel">
		{#if sidePanelView === 'adm'}
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
		{:else if sidePanelView === 'geojson'}
			<h1 class="text-2xl font-bold">GeoJSON Visualizer</h1>
			<div class="flex min-h-0 flex-1 flex-col overflow-x-auto">
				<CodeEditorInput />
			</div>
		{:else}
			Oops! No view selected
		{/if}
	</div>
	<Nav {sidePanelView} onSelectView={(view) => (sidePanelView = view)} />
</div>

<style>
	.side-panel {
		width: 500px;
		flex-shrink: 0;
		display: flex;
		flex-direction: column;
		gap: 10px;
		justify-content: space-between;
	}

	.map {
		box-shadow: 0 0 4px 0 rgba(0, 0, 0, 0.2);
	}
</style>
