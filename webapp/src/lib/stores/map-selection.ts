import { writable } from 'svelte/store';

type SelectableItemId =
	| { type: 'adm-0'; gid0: string }
	| { type: 'adm-1'; gid1: string }
	| { type: 'adm-2'; gid2: string }
	| { type: 'adm-3'; gid3: string }
	| { type: 'adm-4'; gid4: string }
	| { type: 'adm-5'; gid5: string };

type SelectableItemBase<ID extends SelectableItemId, T extends Record<string, unknown>> = ID & T;

type SelectableItem =
	| SelectableItemBase<
			{ type: 'adm-0'; gid0: string },
            { properties: Record<string, unknown> }
	  >
	| SelectableItemBase<
			{ type: 'adm-1'; gid1: string },
            { properties: Record<string, unknown> }
	  >
	| SelectableItemBase<
			{ type: 'adm-2'; gid2: string },
            { properties: Record<string, unknown> }
	  >
	| SelectableItemBase<
			{ type: 'adm-3'; gid3: string },
            { properties: Record<string, unknown> }
	  >
	| SelectableItemBase<
			{ type: 'adm-4'; gid4: string },
            { properties: Record<string, unknown> }
	  >
	| SelectableItemBase<
			{ type: 'adm-5'; gid5: string },
            { properties: Record<string, unknown> }
	  >;

export const mapSelection = writable<SelectableItem | null>(null);
