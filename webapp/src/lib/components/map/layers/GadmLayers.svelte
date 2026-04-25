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
	import GadmLayerHighlight from './GadmLayerHighlight.svelte';
	import { colors } from '$lib/utills/colors';
	import { PUBLIC_MAP_TILES_URL } from '$env/static/public';

	let { map }: { map: maplibregl.Map } = $props();

	let areLayersLoaded = $state(false);

	type LayerInfo = LayerSpecification & {
		url: string;
		source: string;
		id: string;
	};

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
			source: `adm_${level}`,
			'source-layer': `adm_${level}`,
			url: `${PUBLIC_MAP_TILES_URL}/adm_${level}.pmtiles`
		};
		return [
			{
				...baseInfo,
				id: getOutlineLayerIdForAdmLv(level),
				type: 'line',
				paint: {
					'line-color': colors.blackAsh,
					'line-width': [1.5, 1, 0.7, 0.4, 0.2, 0.1][level] ?? 1
				},
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

		areLayersLoaded = true;
	}

	function removeAdmLayers(layerInfos: LayerInfo[]) {
		layerInfos.forEach((info) => {
			removeLayer(map, info.id);
			removeSource(map, info.source);
		});

		areLayersLoaded = false;
	}

	function handleClick(feature: maplibregl.MapGeoJSONFeature) {
		const lv = Number(feature.layer.id.split('-')[1]);
		if (_.isNumber(lv)) {
			mapSelection.set({
				type: 'adm',
				lv,
				featureId: feature.id,
				properties: feature.properties
			});
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
				handleClick(f);
			}
		});
	});

	onDestroy(() => {
		removeAdmLayers(layerInfos);
	});
</script>

{#if areLayersLoaded}
	<GadmLayerHighlight {map} />
{/if}
