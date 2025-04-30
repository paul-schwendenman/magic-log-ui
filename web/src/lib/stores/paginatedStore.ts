import { writable, derived, type Readable } from 'svelte/store';

type Pagination<T> = {
	items: T[];
	page: number;
	pageSize: number;
	totalItems: number;
	totalPages: number;
	hasMore: boolean;
	hasPrev: boolean;
};

export function createPaginationStore<T>(sourceStore: Readable<T[]>, pageSize = 50) {
	const page = writable(0);

	const paginated: Readable<Pagination<T>> = derived(
		[sourceStore, page],
		([$source, $page]) => {
			const start = $page * pageSize;
			const end = start + pageSize;
			const items = $source.slice(start, end);

			return {
				items,
				page: $page,
				pageSize,
				totalItems: $source.length,
				totalPages: Math.ceil($source.length / pageSize),
				hasMore: end < $source.length,
				hasPrev: $page > 0
			};
		},
		{
			items: [],
			page: 0,
			pageSize,
			totalItems: 0,
			totalPages: 0,
			hasMore: false,
			hasPrev: false
		}
	);

	function nextPage() {
		page.update((n) => n + 1);
	}

	function prevPage() {
		page.update((n) => Math.max(0, n - 1));
	}

	function reset() {
		page.set(0);
	}

	return {
		subscribe: paginated.subscribe,
		nextPage,
		prevPage,
		reset
	};
}
