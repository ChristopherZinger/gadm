import { writable } from "svelte/store";

const store = writable<GeoJSON.Feature[] | null>(null);
export  const userProvidedGeometry = {
    subscribe: store.subscribe,
    set: store.set,
    append: (feature: GeoJSON.Feature) => {
        store.update((features) => {
            return [...(features ?? []), feature];
        });
    }
}