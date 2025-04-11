import { writable } from 'svelte/store';
import { browser } from '$app/environment';

export type QueryEntry = {
	query: string;
	ok: boolean;
	timestamp: number;
};

const initial: QueryEntry[] = browser
	? JSON.parse(localStorage.getItem('queryHistory') || '[]')
	: [];

export const queryHistory = writable<QueryEntry[]>(initial);

if (browser) {
	queryHistory.subscribe((value) => {
		localStorage.setItem('queryHistory', JSON.stringify(value));
	});
}

export function addQuery(entry: QueryEntry) {
	queryHistory.update((history) => {
		const filtered = history.filter((item) => item.query !== entry.query);
		return [entry, ...filtered].slice(0, 25);
	});
}

export function clearHistory() {
	queryHistory.set([]);
}
