import { writable, type Writable } from 'svelte/store';
import { browser } from '$app/environment';

export function localStorageStore<T>(key: string, initial: T): Writable<T> {
	const store = writable<T>(initial, (set) => {
		if (browser) {
			const json = localStorage.getItem(key);
			if (json !== null) {
				set(JSON.parse(json));
			}
		}

		return () => {};
	});

	if (browser) {
		store.subscribe((value) => {
			localStorage.setItem(key, JSON.stringify(value));
		});
	}

	return store;
}
