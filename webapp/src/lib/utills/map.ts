import type * as maplibregl from 'maplibre-gl';

export function addSource(map: maplibregl.Map, sourceName: string, url: string) {
	if (!map.getSource(sourceName)) {
		map.addSource(sourceName, {
			type: 'vector',
			url: `pmtiles://${url}`
		});
	}
}

export function removeSource(map: maplibregl.Map, sourceName: string) {
	if (map.getSource(sourceName)) {
		map.removeSource(sourceName);
	}
}

export function addLayer(map: maplibregl.Map, layer: maplibregl.LayerSpecification) {
	if (!map.getLayer(layer.id)) {
		map.addLayer(layer);
	}
}

export function removeLayer(map: maplibregl.Map, layerId: string) {
	if (map.getLayer(layerId)) {
		map.removeLayer(layerId);
	}
}
