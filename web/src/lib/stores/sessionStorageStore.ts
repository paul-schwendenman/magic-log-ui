import { writable, type Writable } from 'svelte/store';
import { browser } from '$app/environment';

export function sessionStorageStore<T>(key: string, initial: T): Writable<T> {
	const start = browser
		? (JSON.parse(sessionStorage.getItem(key) ?? 'null') as T | null) ?? initial
		: initial;

	const store = writable<T>(start);

	if (browser) {
		store.subscribe((value) => {
			sessionStorage.setItem(key, JSON.stringify(value));
		});
	}

	return store;
}
