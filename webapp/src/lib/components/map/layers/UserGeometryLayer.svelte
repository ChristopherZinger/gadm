<script lang="ts">
	import { userProvidedGeometry } from '$lib/stores/user-provided-geometry';
	import {
		USER_GEOMETRY_FILL_LAYER_ID,
		USER_GEOMETRY_OUTLINE_LAYER_ID
	} from '$lib/utills/map-layers-order';
	import type * as maplibregl from 'maplibre-gl';
	import { onMount } from 'svelte';

	let { map }: { map: maplibregl.Map } = $props();

	const isGeojsonSource = (f: maplibregl.Source): f is maplibregl.GeoJSONSource => {
		return f.type === 'geojson';
	};

	const SourceId = 'user-geometry';

	$effect(() => {
		const features = $userProvidedGeometry;
		const src = map.getSource(SourceId);
		if (!src || !isGeojsonSource(src)) {
			return;
		}
		src.setData({
			type: 'FeatureCollection',
			features: features ?? []
		});
	});

	onMount(() => {
		map.addSource(SourceId, {
			type: 'geojson',
			data: {
				type: 'FeatureCollection',
				features: $userProvidedGeometry ?? []
			}
		});

		map.addLayer({
			id: USER_GEOMETRY_FILL_LAYER_ID,
			type: 'fill',
			source: SourceId,
			paint: {
				'fill-color': 'red',
				'fill-opacity': 0.5
			}
		});

		map.addLayer({
			id: USER_GEOMETRY_OUTLINE_LAYER_ID,
			type: 'line',
			source: SourceId,
			paint: {
				'line-color': 'red',
				'line-width': 2
			}
		});

		return () => {
			map.removeSource(SourceId);
			map.removeLayer(USER_GEOMETRY_FILL_LAYER_ID);
			map.removeLayer(USER_GEOMETRY_OUTLINE_LAYER_ID);
		};
	});
</script>
