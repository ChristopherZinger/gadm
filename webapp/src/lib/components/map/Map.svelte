<script lang="ts">
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
			<MapInteraction {map} />
		{/snippet}
	</MapProvider>
</div>
