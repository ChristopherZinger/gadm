<script lang="ts">
	import { basicEditor } from 'prism-code-editor/setups';
	import 'prism-code-editor/prism/languages/markup';
	import { onMount } from 'svelte';
	import _ from 'lodash';
	import { featureCollection } from '@turf/turf';
	import { userProvidedGeometry } from '$lib/stores/user-provided-geometry';
	import { geojsonType } from '@turf/turf';

	const onUpdate = _.debounce((geojsonInput: string) => {
		try {
			const featureCollection = JSON.parse(geojsonInput);
			if (!('features' in featureCollection && Array.isArray(featureCollection.features))) {
				throw new Error('Invalid GeoJSON input');
			}

			featureCollection.features.forEach((feature: GeoJSON.Feature) => {
				geojsonType(feature, 'Feature', 'geojson');
			});

			userProvidedGeometry.set(featureCollection.features);
		} catch (error) {
			console.error(error);
		}
	}, 500);

	onMount(() => {
		basicEditor('#editor', {
			language: 'json',
			theme: 'vs-code-light',
			onUpdate,
			value: JSON.stringify(featureCollection([]), null, 2)
		});
	});
</script>

<div id="editor"></div>
