<script lang="ts">
	import GadmLayers from './layers/GadmLayers.svelte';
	import MapProvider from './MapProvider.svelte';
</script>

<div style="height: 100%; width: 100%;">
	<MapProvider>
		{#snippet children({ map })}
			<button onclick={() => map.setCenter([0, 0])}>Click</button>
			<GadmLayers {map} />
		{/snippet}
	</MapProvider>
</div>
