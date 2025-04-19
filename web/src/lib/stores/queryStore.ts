import { writable, derived, get } from 'svelte/store';

const defaultLimit = 20;

export const query = writable('SELECT * FROM logs ORDER BY timestamp DESC LIMIT 10');
export const page = writable(0);
export const limit = writable(defaultLimit);

export const results = writable<any[]>([]);
export const meta = writable({ hasNextPage: false, hasPreviousPage: false });
export const durationMs = writable<number | null>(null);
export const error = writable<string | null>(null);
export const loading = writable(false);

export async function fetchQuery() {
	error.set(null);
	loading.set(true);
	durationMs.set(null);

	const q = get(query);
	const p = get(page);
	const l = get(limit);

	const start = performance.now();

	try {
		const res = await fetch(`/query?q=${encodeURIComponent(q)}&page=${p}&limit=${l}`);
		if (!res.ok) {
			const text = await res.text();
			throw new Error(text || 'Unknown error');
		}
		const json = await res.json();
		results.set(json.data || []);
		meta.set(json.meta || { hasNextPage: false, hasPreviousPage: false });
		durationMs.set(performance.now() - start);
	} catch (err: any) {
		error.set(err.message || 'Query failed');
	} finally {
		loading.set(false);
	}
}
