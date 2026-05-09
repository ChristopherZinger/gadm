<script lang="ts">
	import type * as maplibregl from 'maplibre-gl';
	import { onDestroy, onMount } from 'svelte';
	import _ from 'lodash';

	type Props = { map: maplibregl.Map; lv: number; featureId: string | null | number };
	let { map, lv, featureId }: Props = $props();

	$effect(() => {
		const id = featureId;
		const timeoutId = setTimeout(() => {
			const filter: maplibregl.FilterSpecification = id ? ['==', ['id'], id] : ['literal', false];

			map.setFilter(`adm-hover-${lv}`, filter);
		}, 150);

		return () => {
			if (timeoutId) {
				clearTimeout(timeoutId);
			}
		};
	});

	onMount(() => {
		const source = `adm_${lv}`;

		map.addLayer({
			id: `adm-hover-${lv}`,
			type: 'line',
			source,
			'source-layer': `adm_${lv}`,
			paint: {
				'line-color': 'red',
				'line-width': 2
			},
			filter: ['literal', false]
		});
	});

	onDestroy(() => {
		map.removeLayer(`adm-hover-${lv}`);
	});
</script>
