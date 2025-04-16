import { derived, get, writable } from 'svelte/store';
import type { LogEntry } from '$lib/types';
import { browser } from '$app/environment';

export const liveLogs = writable<LogEntry[]>([]);
export const liveFilter = writable('');
export const paused = writable(false);
export let buffer = writable<LogEntry[]>([]);

let socket: WebSocket | null = null;
let retryDelay = 1000; // ms
const maxLogs = 500;
let flushTimeout: ReturnType<typeof setTimeout> | null = null;
const flushInterval = 50;

function connect() {
	const ws = new WebSocket(`ws://${location.host}/ws`);
	socket = ws;

	ws.addEventListener('open', () => {
		console.log('✅ WebSocket connected');
		retryDelay = 1000; // Reset backoff
	});

	ws.addEventListener('message', (event) => {
		const entry: LogEntry = JSON.parse(event.data);
		buffer.update(state => [entry, ...state]);

		// Debounce flush
		if (!flushTimeout) {
			flushTimeout = setTimeout(() => {
				if (!get(paused)) {
					liveLogs.update((logs) => {
						const newLogs = [...get(buffer), ...logs];
						buffer.set([]);
						return newLogs.slice(0, maxLogs);
					});
				}
				flushTimeout = null;
			}, flushInterval);
		}
	});

	ws.addEventListener('close', () => {
		console.warn('🔌 WebSocket disconnected — retrying...');
		reconnect();
	});

	ws.addEventListener('error', (err) => {
		console.error('❌ WebSocket error:', err);
		ws.close();
	});
}

function reconnect() {
	setTimeout(() => {
		retryDelay = Math.min(retryDelay * 2, 10000);
		connect();
	}, retryDelay);
}

export const filteredLiveLogs = derived([liveLogs, liveFilter], ([$logs, $filter]) => {
	if (!$filter.trim()) return $logs;

	const f = $filter.toLowerCase();
	return $logs.filter((log) => Object.values(log).join(' ').toLowerCase().includes(f));
});

if (browser) {
	connect();
}
