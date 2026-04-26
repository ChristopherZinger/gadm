<script lang="ts">
	import { onMount } from 'svelte';
	import type * as maplibregl from 'maplibre-gl';
	import { colors } from '$lib/utills/colors';
	import { ADM_LAYER_LEVELS, getOutlineLayerIdForAdmLv } from '$lib/utills/adm-map-layers';

	let {
		map,
		_selection
	}: {
		map: maplibregl.Map;
		_selection: { level: number; featureId: number } | null;
	} = $props();

	const HIGHLIGHT_LAYER_ID = 'adm-highlight';

	function getLayerIdForAdmLv(lv: number): string {
		return `${HIGHLIGHT_LAYER_ID}-${lv}`;
	}

	function handleSelectionChange(selection: { level: number; featureId: number } | null) {
		ADM_LAYER_LEVELS.forEach((lv) => {
			const layerId = getLayerIdForAdmLv(lv);
			if (!map.getLayer(layerId)) {
				// TODO: log error - create logger for prod/local
				return;
			}
			const filter: maplibregl.FilterSpecification =
				selection && lv === selection.level
					? ['==', ['id'], selection.featureId]
					: ['literal', false];

			map.setFilter(layerId, filter);
		});
	}

	$effect(() => {
		handleSelectionChange(_selection);
	});

	onMount(() => {
		ADM_LAYER_LEVELS.forEach((lv) => {
			const source = map.getSource(`adm_${lv}`);
			if (!source) {
				// TODO: log error - create logger for prod/local
				return;
			}
			map.addLayer(
				{
					id: `${HIGHLIGHT_LAYER_ID}-${lv}`,
					source: source.id,
					'source-layer': `adm_${lv}`,
					type: 'fill',
					paint: {
						'fill-color': colors.oceanBlue
					},
					filter: ['literal', false]
				},
				getOutlineLayerIdForAdmLv(5)
			);
		});

		return () => {
			ADM_LAYER_LEVELS.forEach((lv) => {
				if (map.getLayer(getLayerIdForAdmLv(lv))) {
					map.removeLayer(getLayerIdForAdmLv(lv));
				}
			});
		};
	});
</script>
