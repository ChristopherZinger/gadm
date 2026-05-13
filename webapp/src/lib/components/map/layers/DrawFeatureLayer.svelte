<script lang="ts">
	import { drawingModeStore, isGeojsonSource } from '$lib/stores/map-draw-store';
	import { feature } from '@turf/turf';
	import { onMount } from 'svelte';
	import * as maplibregl from 'maplibre-gl';
	import Control from '../controls/MapControl.svelte';
	import { userProvidedGeometry } from '$lib/stores/user-provided-geometry';
	import ShapesIcons from '$lib/icons/ShapesIcons.svelte';
	import { colors } from '$lib/utills/colors';

	let { map }: { map: maplibregl.Map } = $props();

	const SourceId = 'draw-preview';

	function onUpdatePreview(points: GeoJSON.Position[] | null) {
		const src = map.getSource(SourceId);
		if (!src || !isGeojsonSource(src)) {
			return;
		}

		if (!points) {
			src.setData({
				type: 'FeatureCollection',
				features: []
			});
			return;
		}

		src.setData({
			type: 'FeatureCollection',
			features: [
				feature({
					type: 'LineString',
					coordinates: points.map((point) => [point[0], point[1]])
				})
			]
		});
	}
	$effect(() => {
		onUpdatePreview($drawingModeStore?.points ?? null);
	});

	$effect(() => {
		map.getCanvas().style.cursor = $drawingModeStore ? 'crosshair' : 'grab';
	});

	function onMouseMove(e: maplibregl.MapMouseEvent) {
		if ($drawingModeStore) {
			drawingModeStore.setHead([e.lngLat.lng, e.lngLat.lat]);
		}
	}

	function onMouseClick(e: maplibregl.MapMouseEvent) {
		if ($drawingModeStore) {
			drawingModeStore.appendPoint([e.lngLat.lng, e.lngLat.lat]);
		}
	}

	function onDoubleClick() {
		if (!$drawingModeStore) {
			return;
		}
		const points = $drawingModeStore.points;
		userProvidedGeometry.set([
			...($userProvidedGeometry ?? []),
			feature({
				type: 'Polygon',
				coordinates: [[...points, points[0]].map((p) => p)]
			})
		]);
		drawingModeStore.reset();
	}

	onMount(() => {
		if (!map.getSource(SourceId)) {
			map.addSource(SourceId, {
				type: 'geojson',
				data: {
					type: 'FeatureCollection',
					features: []
				}
			});
		}

		if (!map.getLayer('draw-preview_layer')) {
			map.addLayer({
				id: 'draw-preview_layer',
				type: 'line',
				source: SourceId,
				paint: {
					'line-color': 'red',
					'line-width': 2
				}
			});
		}

		map.on('click', onMouseClick);
		map.on('mousemove', onMouseMove);
		map.on('dblclick', onDoubleClick);

		return () => {
			map.off('click', onMouseClick);
			map.off('mousemove', onMouseMove);
			map.off('dblclick', onDoubleClick);
			map.removeLayer('draw-preview');
			map.removeSource('draw-preview');
		};
	});
</script>

<Control {map} position="top-right">
	<button
		onclick={() => {
			drawingModeStore.startDrawingMode('polygon');
		}}
		style="display: flex; align-items: center; justify-content: center;"
	>
		<ShapesIcons height={12} color={colors.blackAsh} />
	</button>
</Control>
