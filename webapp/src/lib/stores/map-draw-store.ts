import { writable } from 'svelte/store';
import type * as maplibregl from 'maplibre-gl';

export type GeometrySketch = GeoJSON.LineString

export const _drawingModeStore = writable<GeometrySketch | null>(null);

export const drawingModeStore = {
	subscribe: _drawingModeStore.subscribe,
	startDrawingMode: function (geometryType: GeometrySketch['type']) {
		_drawingModeStore.set({ type: geometryType, coordinates: [] });
	},
	set: function (d: GeometrySketch) {
		_drawingModeStore.set(d);
	},
	reset: function () {
		_drawingModeStore.set(null);
	},
};

export function isGeojsonSource(f: maplibregl.Source): f is maplibregl.GeoJSONSource {
	return f.type === 'geojson';
}
