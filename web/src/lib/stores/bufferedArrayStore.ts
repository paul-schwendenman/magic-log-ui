import { writable, get } from 'svelte/store';

export function createBufferedLogsStore<T>({
	max = 500,
	flushInterval = 50
} = {}) {
	const logs = writable<T[]>([]);
	let buffer: T[] = [];
	let paused = false;
	let flushTimeout: ReturnType<typeof setTimeout> | null = null;

	function add(entry: T) {
        buffer = [entry, ...buffer].slice(0, max);

		if (!flushTimeout) {
			flushTimeout = setTimeout(() => {
				if (!paused) {
					logs.update((l) => {
						const combined = [...buffer, ...l].slice(0, max);
						buffer = [];
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
		buffer = [];
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
		get bufferSize() {
			return buffer.length;
		},
		get isBufferFull() {
			return buffer.length >= max;
		}
	};
}
