import { browser } from '$app/environment';
import { derived, writable } from 'svelte/store';
import { createBufferedLogsStore } from '$lib/stores/bufferedArrayStore';
import type { LogEntry } from '$lib/types';

export const liveFilter = writable('');
export const paused = writable(false);

export const liveLogs = createBufferedLogsStore<LogEntry>({
	max: 500,
	flushInterval: 100
});

// wire up pause state
paused.subscribe((p) => liveLogs.setPaused(p));

let socket: WebSocket | null = null;

function connect() {
	socket = new WebSocket(`ws://${location.host}/ws`);

	socket.addEventListener('open', () => {
		console.log('✅ WS connected');
	});

	socket.addEventListener('message', (e) => {
		const entry: LogEntry = JSON.parse(e.data);
		liveLogs.add(entry);
	});

	socket.addEventListener('close', () => {
		console.warn('🔌 Disconnected, retrying...');
		setTimeout(connect, 1000);
	});
}

export const filteredLiveLogs = derived(
	[liveLogs, liveFilter],
	([$logs, $filter]) => {
		if (!$filter.trim()) return $logs;
		return $logs.filter((log) =>
			Object.values(log).join(' ').toLowerCase().includes($filter.toLowerCase())
		);
	}
);

if (browser) connect();
