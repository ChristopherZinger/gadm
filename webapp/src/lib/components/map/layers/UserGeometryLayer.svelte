<script lang="ts">
	import { userProvidedGeometry } from '$lib/stores/user-provided-geometry';
	import type * as maplibregl from 'maplibre-gl';
	import { onMount } from 'svelte';

	let { map }: { map: maplibregl.Map } = $props();

	const isGeojsonSource = (f: maplibregl.Source): f is maplibregl.GeoJSONSource => {
		return f.type === 'geojson';
	};

	$effect(() => {
		const features = $userProvidedGeometry;
		const src = map.getSource('user-geometry');
		if (!src || !isGeojsonSource(src)) {
			return;
		}
		src.setData({
			type: 'FeatureCollection',
			features: features ?? []
		});
	});

	onMount(() => {
		map.addSource('user-geometry', {
			type: 'geojson',
			data: {
				type: 'FeatureCollection',
				features: $userProvidedGeometry ?? []
			}
		});

		map.addLayer({
			id: 'user-geometry-fill',
			type: 'fill',
			source: 'user-geometry',
			paint: {
				'fill-color': 'red',
				'fill-opacity': 0.5
			}
		});

		map.addLayer({
			id: 'user-geometry',
			type: 'line',
			source: 'user-geometry',
			paint: {
				'line-color': 'red',
				'line-width': 2
			}
		});

		return () => {
			map.removeSource('user-geometry');
			map.removeLayer('user-geometry-fill');
			map.removeLayer('user-geometry');
		};
	});
</script>
