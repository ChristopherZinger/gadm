<script lang="ts">
	import {
		geometrySketchStore,
		isGeojsonSource,
		convertGeomSketchToPreviewFeature,
		getCompleteFeatureFromGeometrySketch
	} from '$lib/stores/map-draw-store';
	import { onMount } from 'svelte';
	import * as maplibregl from 'maplibre-gl';
	import Control from '../controls/MapControl.svelte';
	import { userProvidedGeometry } from '$lib/stores/user-provided-geometry';
	import ShapesIcons from '$lib/icons/ShapesIcons.svelte';
	import { colors } from '$lib/utills/colors';
	import RectangleIcon from '$lib/icons/RectangleIcon.svelte';
	import CircleIcon from '$lib/icons/CircleIcon.svelte';

	let { map }: { map: maplibregl.Map } = $props();

	const SourceId = 'draw-preview';

	function onUpdatePreview(feature: GeoJSON.Feature<GeoJSON.Geometry> | null) {
		const src = map.getSource(SourceId);
		if (!src || !isGeojsonSource(src)) {
			return;
		}

		src.setData({
			type: 'FeatureCollection',
			features: feature ? [feature] : []
		});
	}

	$effect(() => {
		onUpdatePreview(convertGeomSketchToPreviewFeature($geometrySketchStore));
	});

	$effect(() => {
		map.getCanvas().style.cursor = $geometrySketchStore ? 'crosshair' : 'grab';
	});

	function onMouseMove(e: maplibregl.MapMouseEvent) {
		if (!$geometrySketchStore || $geometrySketchStore.points.length < 1) {
			return;
		}
		const newPoint = [e.lngLat.lng, e.lngLat.lat];
		if ($geometrySketchStore.mode === 'polygon') {
			if ($geometrySketchStore.points.length > 1) {
				geometrySketchStore.setTrailingPoint(newPoint);
			} else {
				geometrySketchStore.appendPoint(newPoint);
			}
			return;
		}

		if ($geometrySketchStore.mode === 'square') {
			if ($geometrySketchStore.points.length > 1) {
				geometrySketchStore.setTrailingPoint(newPoint);
				return;
			}
			geometrySketchStore.appendPoint(newPoint);
			return;
		}

		if ($geometrySketchStore.mode === 'circle') {
			if ($geometrySketchStore.points.length > 1) {
				geometrySketchStore.setTrailingPoint(newPoint);
				return;
			}
			geometrySketchStore.appendPoint(newPoint);
			return;
		}
	}

	function onMouseClick(e: maplibregl.MapMouseEvent) {
		if (!$geometrySketchStore) {
			return;
		}
		const newPoint = [e.lngLat.lng, e.lngLat.lat];
		geometrySketchStore.appendPoint(newPoint);
		onDoneDrawingSquare();
		onDoneDrawingCircle();
	}

	function onDoubleClick() {
		switch ($geometrySketchStore?.mode) {
			case 'square': {
				onDoneDrawingSquare();
				return;
			}
			case 'polygon': {
				onDoneDrawingPolygon();
				return;
			}
			case 'circle': {
				onDoneDrawingCircle();
				return;
			}
		}
	}

	function onDoneDrawingSquare() {
		if ($geometrySketchStore?.mode !== 'square' || $geometrySketchStore.points.length < 2) {
			return;
		}
		const newFeature = getCompleteFeatureFromGeometrySketch($geometrySketchStore);
		userProvidedGeometry.append(newFeature);
		geometrySketchStore.reset();
	}

	function onDoneDrawingCircle() {
		if ($geometrySketchStore?.mode !== 'circle' || $geometrySketchStore.points.length < 2) {
			return;
		}
		const newFeature = getCompleteFeatureFromGeometrySketch($geometrySketchStore);
		userProvidedGeometry.append(newFeature);
		geometrySketchStore.reset();
	}

	function onDoneDrawingPolygon() {
		if ($geometrySketchStore?.mode !== 'polygon') {
			return;
		}
		const newFeature = getCompleteFeatureFromGeometrySketch($geometrySketchStore);
		userProvidedGeometry.append(newFeature);
		geometrySketchStore.reset();
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
			geometrySketchStore.startDrawingMode('polygon');
		}}
		style="display: flex; align-items: center; justify-content: center;"
	>
		<ShapesIcons height={12} color={colors.blackAsh} />
	</button>
	<button
		onclick={() => {
			geometrySketchStore.startDrawingMode('square');
		}}
		style="display: flex; align-items: center; justify-content: center;"
	>
		<RectangleIcon height={12} color={colors.blackAsh} />
	</button>
	<button
		onclick={() => {
			geometrySketchStore.startDrawingMode('circle');
		}}
		style="display: flex; align-items: center; justify-content: center;"
	>
		<CircleIcon height={12} color={colors.blackAsh} />
	</button>
</Control>
