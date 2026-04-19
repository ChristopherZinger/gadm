<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import * as maplibregl from 'maplibre-gl';
	import type { LayerSpecification } from 'maplibre-gl';
	import { addLayer, addSource, removeLayer, removeSource } from '$lib/utills/map';
	import {
		ADM_LAYER_LEVELS,
		getFillLayerIdForAdmLv,
		getOutlineLayerIdForAdmLv
	} from '$lib/utills/adm-map-layers';
	import { LAYERS_IDS_IN_ORDER } from '$lib/utills/map-layers-order';
	import _ from 'lodash';
	import { mapSelection } from '$lib/stores/map-selection';

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

	const layerInfos: LayerInfo[] = ADM_LAYER_LEVELS.map((level) => {
		const baseInfo = {
			source: `ADM_${level}`,
			'source-layer': `ADM_${level}`,
			url: `${PM_TILES_URL}/ADM_${level}.pmtiles`
		};
		return [
			{
				...baseInfo,
				id: getOutlineLayerIdForAdmLv(level),
				type: 'line',
				paint: { 'line-color': 'black' },
				minzoom: admLvToVisibilityZoomLv[level] || 0
			},
			{
				...baseInfo,
				id: getFillLayerIdForAdmLv(level),
				type: 'fill',
				paint: { 'fill-color': 'white' },
				minzoom: admLvToVisibilityZoomLv[level] || 0
			}
		] satisfies LayerInfo[];
	}).flat();

	function createAdmLayers(layerInfos: LayerInfo[]) {
		layerInfos.forEach((info) => {
			addSource(map, info.source, info.url);
		});

		const layersToAdd = [...layerInfos];
		[...LAYERS_IDS_IN_ORDER].reverse().forEach((layerId) => {
			const lInfo = layersToAdd.find((info) => info.id === layerId);
			if (lInfo) {
				addLayer(map, lInfo);
				layersToAdd.splice(layersToAdd.indexOf(lInfo), 1);
			}
		});

		layersToAdd.forEach((info) => {
			addLayer(map, info);
		});
	}

	function removeAdmLayers(layerInfos: LayerInfo[]) {
		layerInfos.forEach((info) => {
			removeLayer(map, info.id);
			removeSource(map, info.source);
		});
	}

	function handleClick(feature: maplibregl.MapGeoJSONFeature) {
		const lv = feature.layer.id.split('-')[1];

		switch (lv) {
			case '0':
				mapSelection.set({
					type: 'adm-0',
					gid0: feature.properties.gid_0,
					properties: feature.properties
				});
				break;
			case '1':
				mapSelection.set({
					type: 'adm-1',
					gid1: feature.properties.gid_1,
					properties: feature.properties
				});
				break;
			case '2':
				mapSelection.set({
					type: 'adm-2',
					gid2: feature.properties.gid_2,
					properties: feature.properties
				});
				break;
			case '3':
				mapSelection.set({
					type: 'adm-3',
					gid3: feature.properties.gid_3,
					properties: feature.properties
				});
				break;
			case '4':
				mapSelection.set({
					type: 'adm-4',
					gid4: feature.properties.gid_4,
					properties: feature.properties
				});
				break;
			case '5':
				mapSelection.set({
					type: 'adm-5',
					gid5: feature.properties.gid_5,
					properties: feature.properties
				});
				break;
			default:
				mapSelection.set(null);
		}
	}

	onMount(() => {
		createAdmLayers(layerInfos);

		map.on('click', (e) => {
			const features = map.queryRenderedFeatures(e.point, {
				layers: ADM_LAYER_LEVELS.map(getFillLayerIdForAdmLv)
			});
			const f = _.sortBy(features, (f) => f.layer.id.split('-')[1]).reverse()[0];
			if (f) {
				console.log(f);
				handleClick(f);
			}
		});
	});

	onDestroy(() => {
		removeAdmLayers(layerInfos);
	});
</script>
