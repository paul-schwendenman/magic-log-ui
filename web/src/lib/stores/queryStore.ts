import type { TimeRange } from '$lib/types';
import { writable, derived, type Writable, type Readable } from 'svelte/store';

interface QueryResult<T> {
	error: string | null;
	results: T[];
	meta: {
		hasNextPage: boolean;
		hasPreviousPage: boolean;
		page: number;
		totalPages: number;
	};
	durationMs: number | null;
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
	let meta = { hasNextPage: false, hasPreviousPage: false, page: 0, totalPages: 1 };
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

	const timeRangeStore = writable<TimeRange>({
		from: new Date(Date.now() - 15 * 60 * 1000),
		to: new Date()
	});

	queryStore.subscribe(() => pageStore.set(0));
	limitStore.subscribe(() => pageStore.set(0));
	timeRangeStore.subscribe(() => pageStore.set(0));

	const results = derived<
		[Readable<number>, Readable<number>, Readable<string>, Readable<TimeRange>],
		QueryResult<any>
	>(
		[pageStore, limitStore, queryStore, timeRangeStore],
		([$pageStore, $limitStore, $queryStore, $timeRangeStore], set) => {
			const query = new URLSearchParams({
				q: $queryStore,
				limit: $limitStore.toString(),
				page: $pageStore.toString(),
				from: $timeRangeStore.from.toISOString(),
				to: $timeRangeStore.to.toISOString()
			});

			fetchQuery(`/query?${query.toString()}`, loading).then((r) => set({ ...r }));
		},
		{
			error: null,
			results: [],
			meta: { hasNextPage: false, hasPreviousPage: false, page: 0, totalPages: 1 },
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
		setTimeRange: timeRangeStore.set,
		page: pageStore,
		limit: limitStore,
		nextPage,
		prevPage
	};
};
