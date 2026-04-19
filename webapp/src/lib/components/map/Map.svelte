<script lang="ts">
	import GadmLayers from './layers/GadmLayers.svelte';
	import MapProvider from './MapProvider.svelte';
</script>

<div style="height: 100%; width: 100%;">
	<MapProvider>
		{#snippet children({ map })}
			<GadmLayers {map} />
		{/snippet}
	</MapProvider>
</div>
