export function appendPointToGeometry(geometry: GeoJSON.LineString, position: GeoJSON.Position) {
	return {
		...geometry,
		coordinates: [...geometry.coordinates, position]
	};
}

export function setGeometryTrailingPoint(geometry: GeoJSON.LineString, position: GeoJSON.Position) {
	return {
		...geometry,
		coordinates: [...geometry.coordinates.slice(0, -1), position]
	};
}

export function convertLineStringToPolygon(geometry: GeoJSON.LineString): GeoJSON.Polygon {
	return {
		type: 'Polygon',
		coordinates: [[...geometry.coordinates, geometry.coordinates[0]]]
	};
}
