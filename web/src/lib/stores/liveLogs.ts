import { writable } from 'svelte/store';
import type { LogEntry } from '$lib/types';
import { browser } from '$app/environment';

export const liveLogs = writable<LogEntry[]>([]);

let socket: WebSocket | null = null;
let retryDelay = 1000; // ms
const maxLogs = 500;
let buffer: LogEntry[] = [];
let flushTimeout: ReturnType<typeof setTimeout> | null = null;

function connect() {
	const ws = new WebSocket(`ws://${location.host}/ws`);
	socket = ws;

	ws.addEventListener('open', () => {
		console.log('✅ WebSocket connected');
		retryDelay = 1000; // Reset backoff
	});

	ws.addEventListener('message', (event) => {
		const entry: LogEntry = JSON.parse(event.data);
		buffer.unshift(entry);

		// Debounce flush
		if (!flushTimeout) {
			flushTimeout = setTimeout(() => {
				liveLogs.update((logs) => {
					const newLogs = [...buffer, ...logs];
					buffer = [];
					flushTimeout = null;
					return newLogs.slice(0, maxLogs);
				});
			}, 50); // adjust for feel
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

if(browser) {
    connect();
}
