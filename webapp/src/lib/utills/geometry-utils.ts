export function convertTwoPointsToSquarePolygon(
	points: GeoJSON.Position[]
): GeoJSON.Polygon | null {
	if (points.length < 2) {
		console.warn('Requires at least two points to convert to a square polygon');
		return null;
	}
	const [x1, y1] = points[0];
	const [x2, y2] = points[1];
	const square: GeoJSON.Polygon = {
		type: 'Polygon',
		coordinates: [
			[
				[x1, y1],
				[x2, y1],
				[x2, y2],
				[x1, y2],
				[x1, y1]
			]
		]
	};
	return square;
}
