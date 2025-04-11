import { writable } from 'svelte/store';
import { browser } from '$app/environment';

const initial = browser
	? JSON.parse(localStorage.getItem('queryHistory') || '[]')
	: [];

export const queryHistory = writable<string[]>(initial);

if (browser) {
	queryHistory.subscribe((value) => {
		localStorage.setItem('queryHistory', JSON.stringify(value));
	});
}

export function addQuery(q: string) {
	queryHistory.update((history) => {
		if (!q.trim()) return history;
		const deduped = history.filter((item) => item !== q);
		return [q, ...deduped].slice(0, 25); // max 25
	});
}

export function clearHistory() {
	queryHistory.set([]);
}
