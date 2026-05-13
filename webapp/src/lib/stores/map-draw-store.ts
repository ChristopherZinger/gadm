import { writable } from 'svelte/store';
import type * as maplibregl from 'maplibre-gl';

export const _drawingModeStore = writable<{
	mode: 'polygon';
	points: GeoJSON.Position[];
} | null>(null);

export const drawingModeStore = {
	subscribe: _drawingModeStore.subscribe,
	startDrawingMode: function (mode: 'polygon') {
		_drawingModeStore.set({ mode, points: [] });
	},
	appendPoint: (point: GeoJSON.Position) => {
		_drawingModeStore.update((values) => {
			if (!values) {
				return values;
			}
			return { ...values, points: [...values.points, point] };
		});
	},
	reset: function () {
		_drawingModeStore.set(null);
	},
	setHead: function (point: GeoJSON.Position) {
		_drawingModeStore.update((values) => {
			if (!values) {
				return values;
			}
			return { ...values, points: [...values.points.slice(0, -1), point] };
		});
	}
};

export function isGeojsonSource(f: maplibregl.Source): f is maplibregl.GeoJSONSource {
	return f.type === 'geojson';
}
