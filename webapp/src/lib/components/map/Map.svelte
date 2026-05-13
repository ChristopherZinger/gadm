<script lang="ts">
	import 'maplibre-gl/dist/maplibre-gl.css';
	import DrawFeatureLayer from './layers/DrawFeatureLayer.svelte';
	import GadmLayers from './layers/GadmLayers.svelte';
	import UserGeometryLayer from './layers/UserGeometryLayer.svelte';
	import MapProvider from './MapProvider.svelte';
	import MapInteraction from './MapInteraction.svelte';
</script>

<div style="height: 100%; width: 100%;">
	<MapProvider>
		{#snippet children({ map })}
			<GadmLayers {map} />
			<UserGeometryLayer {map} />
			<DrawFeatureLayer {map} />
			<MapInteraction {map} />
		{/snippet}
	</MapProvider>
</div>
