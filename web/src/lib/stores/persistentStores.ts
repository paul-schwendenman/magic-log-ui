import { writable, type Writable } from 'svelte/store';
import { browser } from '$app/environment';

export function persistentStore<T>(
	key: string,
	initial: T,
	storage: Storage
): Writable<T> {
	const start = browser
		? (JSON.parse(storage.getItem(key) ?? 'null') as T | null) ?? initial
		: initial;

	const store = writable<T>(start);

	if (browser) {
		store.subscribe((value) => {
			storage.setItem(key, JSON.stringify(value));
		});
	}

	return store;
}
