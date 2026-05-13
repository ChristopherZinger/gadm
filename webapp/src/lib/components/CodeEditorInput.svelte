<script lang="ts">
	import { basicEditor } from 'prism-code-editor/setups';
	import 'prism-code-editor/prism/languages/markup';
	import { onMount } from 'svelte';
	import _ from 'lodash';
	import { featureCollection } from '@turf/turf';
	import { userProvidedGeometry } from '$lib/stores/user-provided-geometry';
	import { geojsonType } from '@turf/turf';
	import type { PrismEditor } from 'prism-code-editor';

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

	let editor: PrismEditor | undefined;
	$effect(() => {
		const features = $userProvidedGeometry ?? [];
		if (editor) {
			editor.setOptions({
				value: JSON.stringify(featureCollection(features), null, 2)
			});
		}
	});

	onMount(() => {
		const features = $userProvidedGeometry ?? [];
		editor = basicEditor('#editor', {
			language: 'json',
			theme: 'vs-code-light',
			onUpdate,
			value: JSON.stringify(featureCollection(features), null, 2)
		});
	});
</script>

<div id="editor"></div>
