import { writable } from "svelte/store";

export const userProvidedGeometry = writable<GeoJSON.Feature[] | null>(null);