<script lang="ts">
	import { onMount } from 'svelte';
	import * as maplibregl from 'maplibre-gl';
	import { Protocol } from 'pmtiles';

	onMount(() => {
		const protocol = new Protocol();
		maplibregl.addProtocol('pmtiles', protocol.tile);

		const map = new maplibregl.Map({
			container: 'map',
			style: 'https://demotiles.maplibre.org/style.json',
			center: [0, 0],
			zoom: 2
		});

		const sourceName = 'pmt-adm0';

		map.on('load', () => {
			const tileHttpUrl = `http://127.0.0.1:8080/adm_0.pmtiles`;
			map.addSource(sourceName, {
				type: 'vector',
				url: `pmtiles://${tileHttpUrl}`
			});

			// `source-layer` must match a layer name inside the MVT (see `npx pmtiles show adm_0.pmtiles`).
			map.addLayer({
				id: 'pmtiles-adm0-lines',
				type: 'line',
				source: sourceName,
				'source-layer': 'zcta',
				paint: {
					'line-color': 'red'
				}
			});

			map.addLayer({
				id: 'pmtiles-adm0-fill',
				type: 'fill',
				source: sourceName,
				'source-layer': 'zcta',
				paint: {
					'fill-opacity': 0.1,
					'fill-color': 'blue'
				}
			});

			map.getStyle().layers.forEach((layer) => {
				console.log(layer);
				if (!['pmtiles-adm0-fill', 'pmtiles-adm0-lines'].includes(layer.id)) {
					map.setLayoutProperty(layer.id, 'visibility', 'none');
				}
			});
		});

		return () => {
			map.remove();
			maplibregl.removeProtocol('pmtiles');
		};
	});
</script>

<div id="map" style="height: 100vh; width: 90vw;"></div>
