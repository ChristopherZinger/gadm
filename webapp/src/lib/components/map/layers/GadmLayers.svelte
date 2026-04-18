<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import * as maplibregl from 'maplibre-gl';
	import type { LayerSpecification } from 'maplibre-gl';
	import { addLayer, addSource, removeLayer, removeSource } from '$lib/utills/map';

	let { map }: { map: maplibregl.Map } = $props();

	type LayerInfo = LayerSpecification & {
		url: string;
		source: string;
		id: string;
	};

	const PM_TILES_URL = 'http://localhost:8080/pmtiles';

	const admLvToVisibilityZoomLv: Record<number, number> = {
		0: 0,
		1: 2,
		2: 4,
		3: 6,
		4: 8,
		5: 10
	};

	const layerInfos: LayerInfo[] = [0, 1, 2, 3, 4, 5].reverse().map((level) => {
		const baseInfo = {
			source: `ADM_${level}`,
			'source-layer': `ADM_${level}`,
			url: `${PM_TILES_URL}/ADM_${level}.pmtiles`
		};
		return {
			...baseInfo,
			id: `adm-${level}-line`,
			type: 'line',
			paint: {
				'line-color': 'black'
			},
			minzoom: admLvToVisibilityZoomLv[level] || 0
		};
	});

	function createAdmLayers(layerInfos: LayerInfo[]) {
		layerInfos.forEach((info) => {
			addSource(map, info.source, info.url);
			addLayer(map, info);
		});
	}

	function removeAdmLayers(layerInfos: LayerInfo[]) {
		layerInfos.forEach((info) => {
			removeSource(map, info.source);
			removeLayer(map, info.id);
		});
	}

	function addAdmBgFillLayer() {
		addLayer(map, {
			source: `ADM_0`,
			'source-layer': `ADM_0`,
			id: 'adm-0-fill',
			type: 'fill',
			paint: { 'fill-color': 'black', 'fill-opacity': 0.1 }
		});
	}

	onMount(() => {
		createAdmLayers(layerInfos);

		addAdmBgFillLayer();
	});

	onDestroy(() => {
		removeAdmLayers(layerInfos);
	});
</script>
