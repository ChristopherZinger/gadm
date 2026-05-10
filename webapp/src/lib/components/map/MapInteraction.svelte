<script lang="ts">
	import { ADM_LAYER_LEVELS, getFillLayerIdForAdmLv } from '$lib/utills/adm-map-layers';
	import { USER_GEOMETRY_FILL_LAYER_ID } from '$lib/utills/map-layers-order';
	import type * as maplibregl from 'maplibre-gl';
	import { onDestroy, onMount } from 'svelte';
	import { hoveredFeature, selectedAdmInfo, selectedFeature } from '$lib/stores/map-selection';

	let { map }: { map: maplibregl.Map } = $props();

	const onMouseClick = (e: maplibregl.MapMouseEvent) => {
		const mapFeatures = map.queryRenderedFeatures(e.point, {
			layers: getEligibleInteractiveLayersInOrder()
		});
		const feature = findMostRelevantFeature(mapFeatures);
		selectedFeature.set(feature ? { featureId: feature.id, layerId: feature.layer.id } : null);
		// deprecate this and retrieve from api
		selectedAdmInfo.set(feature ? feature.properties : null);
	};

	const onMouseMove = (e: maplibregl.MapMouseEvent) => {
		const mapFeatures = map.queryRenderedFeatures(e.point, {
			layers: getEligibleInteractiveLayersInOrder()
		});
		const feature = findMostRelevantFeature(mapFeatures);
		hoveredFeature.set(feature ? { featureId: feature.id, layerId: feature.layer.id } : null);
	};

	function findMostRelevantFeature(
		mapFeatures: maplibregl.MapGeoJSONFeature[]
	): maplibregl.MapGeoJSONFeature | null {
		for (const interactiveLayer of getEligibleInteractiveLayersInOrder()) {
			const f = mapFeatures.find((f) => f.layer.id === interactiveLayer);
			if (f) {
				return f;
			}
		}
		return null;
	}

	function getEligibleInteractiveLayersInOrder() {
		const result = [...InteractiveLayersInOrder].reverse().filter((lId) => {
			const l = map.getLayer(lId);
			return l && l.visibility !== 'none';
		});
		return result;
	}

	const InteractiveLayersInOrder = [
		...ADM_LAYER_LEVELS.map(getFillLayerIdForAdmLv),
		USER_GEOMETRY_FILL_LAYER_ID
	];

	onMount(() => {
		map.on('mousemove', onMouseMove);
		map.on('click', onMouseClick);
	});

	onDestroy(() => {
		map.off('mousemove', onMouseMove);
		map.off('click', onMouseClick);
	});
</script>
