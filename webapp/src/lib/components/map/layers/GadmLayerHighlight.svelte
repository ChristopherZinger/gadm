<script lang="ts">
	import { get } from 'svelte/store';
	import { mapSelection, type SelectableItem } from '$lib/stores/map-selection';
	import { onMount } from 'svelte';
	import type * as maplibregl from 'maplibre-gl';
	import _ from 'lodash';
	import { colors } from '$lib/utills/colors';
	import { getOutlineLayerIdForAdmLv } from '$lib/utills/adm-map-layers';

	let { map }: { map: maplibregl.Map } = $props();

	const HIGHLIGHT_LAYER_ID = 'adm-highlight';

	function applyHighlightFilter(selection?: SelectableItem | null) {
		const layerId = `${HIGHLIGHT_LAYER_ID}-${selection?.lv}`;
		if (!map.getLayer(layerId) || selection?.type !== 'adm') {
			return;
		}
		const _id = Number(selection?.featureId);

		const filter: maplibregl.FilterSpecification = _.isNumber(_id)
			? ['==', ['id'], _id]
			: ['literal', false];
		map.setFilter(layerId, filter);
	}

	$effect(() => {
		applyHighlightFilter($mapSelection);
	});

	onMount(() => {
		_.range(0, 6).forEach((lv) => {
			map.addLayer(
				{
					id: `${HIGHLIGHT_LAYER_ID}-${lv}`,
					source: `adm_${lv}`,
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
	});
</script>
