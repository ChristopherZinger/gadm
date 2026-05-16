import { writable } from 'svelte/store';
import type * as maplibregl from 'maplibre-gl';
import { feature, circle, distance } from '@turf/turf';
import { convertTwoPointsToSquarePolygon } from '$lib/utills/geometry-utils';

export type GeometrySketch = {
	mode: 'square' | 'polygon' | 'circle';
	points: GeoJSON.Position[];
};

export const _store = writable<GeometrySketch | null>(null);

export const geometrySketchStore = {
	subscribe: _store.subscribe,
	startDrawingMode: function (sketchType: GeometrySketch['mode']): null {
		switch (sketchType) {
			case 'polygon':
				_store.set({ mode: 'polygon', points: [] });
				return null;
			case 'square':
				_store.set({ mode: 'square', points: [] });
				return null;
			case 'circle':
				_store.set({ mode: 'circle', points: [] });
				return null;
		}
	},
	appendPoint: function (point: GeoJSON.Position) {
		_store.update((v) => {
			if (!v) {
				return null;
			}
			if ((v.mode === 'square' || v.mode === 'circle') && v.points.length > 1) {
				return v;
			}
			return { ...v, points: [...v.points, point] };
		});
	},
	setTrailingPoint: function (point: GeoJSON.Position) {
		_store.update((v) => {
			if (!v) {
				return null;
			}
			return { ...v, points: [...v.points.slice(0, -1), point] };
		});
	},
	reset: function () {
		_store.set(null);
	}
};

export function isGeojsonSource(f: maplibregl.Source): f is maplibregl.GeoJSONSource {
	return f.type === 'geojson';
}

export function convertGeomSketchToPreviewFeature(
	sketch: GeometrySketch | null
): GeoJSON.Feature<GeoJSON.Geometry> | null {
	if (!sketch) {
		return null;
	}
	switch (sketch.mode) {
		case 'polygon':
			return feature({
				type: 'LineString',
				coordinates: sketch.points
			});
		case 'square':
			return feature(convertTwoPointsToSquarePolygon(sketch.points));
		case 'circle': {
			if (sketch.points.length < 2) {
				return null
			}
			const units = 'meters' as const;
			return circle(sketch.points[0], distance(sketch.points[0], sketch.points[1], { units }), {
				steps: 64,
				units
			});
		}
	}
}

export function getCompleteFeatureFromGeometrySketch(
	sketch: GeometrySketch
): GeoJSON.Feature<GeoJSON.Geometry> {
	switch (sketch.mode) {
		case 'polygon':
			return feature({
				type: 'Polygon',
				coordinates: [[...sketch.points, sketch.points[0]]]
			});
		case 'square':
			return feature(convertTwoPointsToSquarePolygon(sketch.points));
		case 'circle': {
			const units = 'meters' as const;
			return circle(sketch.points[0], distance(sketch.points[0], sketch.points[1], { units }), {
				steps: 64,
				units
			});
		}
	}
}
