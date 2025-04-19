import { number } from '$lib/paraglide/registry';
import { writable, derived, type Writable, type Readable } from 'svelte/store';

interface QueryResult<T> {
	error: string | null;
	results: T[];
	meta: { hasNextPage: boolean; hasPreviousPage: boolean };
	durationMs: number | null;
	page: number;
}

const defaultLimit = 20;

export const loading = writable(false);

export async function fetchQuery<T>(
	url: string,
	loading: Writable<boolean> = writable(false)
): Promise<QueryResult<T>> {
	let error = null;
	let durationMs = null;
	let results: T[] = [];
	let meta = { hasNextPage: false, hasPreviousPage: false };
	loading.set(true);

	const start = performance.now();

	try {
		const res = await fetch(url);
		if (!res.ok) {
			const text = await res.text();
			throw new Error(text || 'Unknown error');
		}
		const json = await res.json();
		results = json.data || [];
		meta = json.meta || { hasNextPage: false, hasPreviousPage: false };

		durationMs = performance.now() - start;
	} catch (err: any) {
		error = err.message || 'Query failed';
	} finally {
		loading.set(false);

		return {
			error,
			results,
			meta,
			durationMs
		};
	}
}

interface QueryStoreProps {
	query: string;
	page?: number;
	limit?: number;
}

export const createQueryStore = ({
	query: initialQuery,
	page: initialPage = 0,
	limit: initialLimit = defaultLimit
}: QueryStoreProps) => {
	const pageStore = writable(initialPage);
	const limitStore = writable(initialLimit);
	const queryStore = writable(initialQuery);
	const loading = writable(false);

	queryStore.subscribe(() => pageStore.set(0));
	limitStore.subscribe(() => pageStore.set(0));

	const results = derived<[Readable<number>, Readable<number>, Readable<string>], QueryResult<any>>(
		[pageStore, limitStore, queryStore],
		([$pageStore, $limitStore, $queryStore], set) => {
			const query = new URLSearchParams({
				q: $queryStore,
				limit: $limitStore.toString(),
				page: $pageStore.toString()
			});

			fetchQuery(`/query?${query.toString()}`, loading).then((r) =>
				set({ ...r, page: $pageStore })
			);
		},
		{
			error: null,
			results: [],
			meta: { hasNextPage: false, hasPreviousPage: false },
			durationMs: null
		}
	);

	function nextPage() {
		pageStore.update((n) => n + 1);
	}
	function prevPage() {
		pageStore.update((n) => Math.max(0, n - 1));
	}

	return {
		subscribe: results.subscribe,
		setQuery: queryStore.set,
		setLimit: limitStore.set,
		page: pageStore,
		limit: limitStore,
		nextPage,
		prevPage
	};
};
