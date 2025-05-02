import { localStorageStore } from './stores/localStorageStore';

export type QueryEntry = {
	query: string;
	ok: boolean;
	timestamp: number;
};

const initial: QueryEntry[] = [];

export const queryHistory = localStorageStore<QueryEntry[]>('queryHistory', initial);

export function addQuery(entry: QueryEntry) {
	queryHistory.update((history) => {
		const filtered = history.filter((item) => item.query !== entry.query);
		return [entry, ...filtered].slice(0, 25);
	});
}

export function clearHistory() {
	queryHistory.set([]);
}
