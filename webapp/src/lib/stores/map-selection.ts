import { writable } from 'svelte/store';

export const selectedAdmInfo = writable<Record<string, unknown> | null>(null);

type SelectedFeature = {
	featureId: string | undefined | number;
	layerId: string;
};
export const selectedFeature = writable<SelectedFeature | null>(null);

type HoverFeature = {
	featureId: string | undefined | number;
	layerId: string;
};
export const hoveredFeature = writable<HoverFeature | null>(null);