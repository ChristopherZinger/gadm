import { writable } from 'svelte/store';

type SelectableAdmId = {
	type: 'adm';
	featureId: string | undefined | number;
	lv: number;
};

export type SelectableItem = SelectableAdmId & { properties: Record<string, unknown> };

export const mapSelection = writable<SelectableItem | null>(null);

type HoverFeature = {
	featureId: string | undefined | number;
	layerId: string;
};

export const hoveredFeature = writable<HoverFeature | null>(null);