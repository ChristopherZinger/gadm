<script lang="ts">
	import type * as maplibregl from 'maplibre-gl';
	import { colors } from '$lib/utills/colors';
	import { getOutlineLayerIdForAdmLv } from '$lib/utills/adm-map-layers';

	let {
		map,
		featureId,
		lv
	}: {
		map: maplibregl.Map;
		lv: number;
		featureId: string | null | number;
	} = $props();

	const HIGHLIGHT_LAYER_ID_PREFIX = 'adm-highlight';
	function getLayerIdForAdmLv(lv: number): string {
		return `${HIGHLIGHT_LAYER_ID_PREFIX}-${lv}`;
	}
	let layerId = $derived(getLayerIdForAdmLv(lv));

	// should this also take layerId ?
	function handleSelectionChange(featureId: string | number | null) {
		const filter: maplibregl.FilterSpecification = featureId
			? ['==', ['id'], featureId]
			: ['literal', false];

		map.setFilter(layerId, filter);
	}

	$effect(() => {
		const source = map.getSource(`adm_${lv}`);
		if (!source) {
			// TODO: log error - create logger for prod/local
			return;
		}
		map.addLayer(
			{
				id: layerId,
				source: source.id,
				'source-layer': `adm_${lv}`,
				type: 'fill',
				paint: {
					'fill-color': colors.col4
				},
				filter: ['literal', false]
			},
			getOutlineLayerIdForAdmLv(5)
		);

		return () => {
			if (map.getLayer(layerId)) {
				map.removeLayer(layerId);
			}
		};
	});

	$effect(() => {
		handleSelectionChange(featureId ? Number(featureId) : null);
	});
</script>
