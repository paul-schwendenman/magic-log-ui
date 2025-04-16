import { writable, get, derived, readable } from 'svelte/store';

export function createBufferedLogsStore<T>({ max = 500, flushInterval = 50 } = {}) {
	const logs = writable<T[]>([]);
	const buffer = writable<T[]>([]);
	const bufferSize = derived([buffer], ([b]) => b.length);

	let paused = false;
	let flushTimeout: ReturnType<typeof setTimeout> | null = null;

	function add(entry: T) {
		buffer.update((b) => [entry, ...b].slice(0, max));

		if (!flushTimeout) {
			flushTimeout = setTimeout(() => {
				if (!paused) {
					logs.update((l) => {
						const combined = [...get(buffer), ...l].slice(0, max);
						buffer.set([]);
						return combined;
					});
				}
				flushTimeout = null;
			}, flushInterval);
		}
	}

	function clearLogs() {
		logs.set([]);
	}

	function clearBuffer() {
		buffer.set([]);
	}

	function setPaused(p: boolean) {
		paused = p;
	}

	return {
		subscribe: logs.subscribe,
		add,
		clearLogs,
		clearBuffer,
		setPaused,
		bufferSize,
		isBufferFull: derived(bufferSize, (bufferSize) => bufferSize >= max)
	};
}
