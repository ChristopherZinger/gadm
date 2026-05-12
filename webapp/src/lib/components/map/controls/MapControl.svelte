<script lang="ts">
	import MapControlContainer from './MapControlContainer.svelte';
	import type { Snippet } from 'svelte';
	import { mount, onDestroy, onMount, unmount } from 'svelte';

	type Props = {
		position: maplibregl.ControlPosition;
		map: maplibregl.Map;
		children: Snippet;
	};
	let { position, map, children }: Props = $props();

	class Control implements maplibregl.IControl {
		container!: HTMLDivElement;
		app?: Record<string, unknown>;
		onAdd() {
			this.container = document.createElement('div');
			this.container.className = 'maplibregl-ctrl maplibregl-ctrl-group';
			this.app = mount(MapControlContainer, {
				target: this.container,
				props: { content: controlUi }
			});
			return this.container;
		}
		onRemove() {
			if (this.app) unmount(this.app);
			this.container.remove();
		}
	}

	onMount(() => {
		map.addControl(new Control(), position);
	});

	onDestroy(() => {
		map.removeControl(new Control());
	});
</script>

{#snippet controlUi()}
	{@render children()}
{/snippet}
