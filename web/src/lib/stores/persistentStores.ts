import { writable, type Writable } from 'svelte/store';
import { browser } from '$app/environment';

export function persistentStore<T>(
	key: string,
	initial: T,
	getStorage: () => Storage
): Writable<T> {
	const start = browser
		? ((JSON.parse(getStorage().getItem(key) ?? 'null') as T | null) ?? initial)
		: initial;

	const store = writable<T>(start);

	if (browser) {
		const storage = getStorage();

		store.subscribe((value) => {
			storage.setItem(key, JSON.stringify(value));
		});
	}

	return store;
}
