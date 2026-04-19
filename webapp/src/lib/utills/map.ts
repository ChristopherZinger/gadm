import type * as maplibregl from 'maplibre-gl';
import { LAYERS_IDS_IN_ORDER } from './map-layers-order';

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
	if (map.getLayer(layer.id)) {
        return;
	}
    const beforeId = getLayerIdAbove(layer.id)
	map.addLayer(layer, beforeId);
}

export function removeLayer(map: maplibregl.Map, layerId: string) {
	if (map.getLayer(layerId)) {
		map.removeLayer(layerId);
	}
}


function getLayerIdAbove (layerId: string) {
    const index = LAYERS_IDS_IN_ORDER.indexOf(layerId);
    if (index === -1) {
        return undefined;
    }
    return LAYERS_IDS_IN_ORDER[index + 1] ?? undefined;
}