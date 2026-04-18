<script lang="ts">
	import type { Snippet } from 'svelte';
	import { onMount } from 'svelte';
	import * as maplibregl from 'maplibre-gl';
	import type { Map } from 'maplibre-gl';
	import { Protocol } from 'pmtiles';

	let {
		children
	}: {
		children: Snippet<[{ map: Map }]>;
	} = $props();

	let map = $state<maplibregl.Map | null | undefined>(undefined);

	onMount(() => {
		const protocol = new Protocol();
		maplibregl.addProtocol('pmtiles', protocol.tile);

		const _map = new maplibregl.Map({
			container: 'map',
			style: baseStyle,
			center: [0, 0],
			zoom: 2
		});

		_map.on('load', () => {
			map = _map;
		});

		return () => {
			map = null;
			_map.remove();
			maplibregl.removeProtocol('pmtiles');
		};
	});

	const baseStyle = {
		id: '43f36e14-e3f5-43c1-84c0-50a9c80dc5c7',
		name: 'MapLibre',
		zoom: 0.8619833357855968,
		pitch: 0,
		center: [17.65431710431244, 32.954120326746775],
		glyphs: 'https://demotiles.maplibre.org/font/{fontstack}/{range}.pbf',
		layers: [
			{
				id: 'background',
				type: 'background',
				paint: {
					'background-color': 'white'
				},
				filter: ['all'],
				layout: {
					visibility: 'visible'
				},
				maxzoom: 24
			},
			{
				id: 'geolines',
				type: 'line',
				paint: {
					'line-color': '#1077B0',
					'line-opacity': 1,
					'line-dasharray': [3, 3]
				},
				filter: ['all', ['!=', 'name', 'International Date Line']],
				layout: {
					visibility: 'visible'
				},
				source: 'maplibre',
				maxzoom: 24,
				'source-layer': 'geolines'
			},
			{
				id: 'geolines-label',
				type: 'symbol',
				paint: {
					'text-color': '#1077B0',
					'text-halo-blur': 1,
					'text-halo-color': 'rgba(255, 255, 255, 1)',
					'text-halo-width': 1
				},
				filter: ['all', ['!=', 'name', 'International Date Line']],
				layout: {
					'text-font': ['Open Sans Semibold'],
					'text-size': {
						stops: [
							[2, 12],
							[6, 16]
						]
					},
					'text-field': '{name}',
					visibility: 'visible',
					'symbol-placement': 'line'
				},
				source: 'maplibre',
				maxzoom: 24,
				minzoom: 1,
				'source-layer': 'geolines'
			}
		],
		bearing: 0,
		sources: {
			maplibre: {
				url: 'https://demotiles.maplibre.org/tiles/tiles.json',
				type: 'vector'
			}
		},
		version: 8,
		metadata: {
			'maptiler:copyright':
				'This style was generated on MapTiler Cloud. Usage is governed by the license terms in https://github.com/maplibre/demotiles/blob/gh-pages/LICENSE',
			'openmaptiles:version': '3.x'
		}
	};
</script>

<div id="map" style="height: 100%; width: 100%;">
	{#if map === undefined}
		<div>Loading map...</div>
	{:else if map === null}
		<div>Error loading map</div>
	{:else}
		{@render children({ map })}
	{/if}
</div>
