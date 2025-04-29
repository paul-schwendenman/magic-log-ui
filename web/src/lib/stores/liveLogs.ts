import { browser } from '$app/environment';
import { derived, writable } from 'svelte/store';
import { createBufferedLogsStore } from '$lib/stores/bufferedArrayStore';
import type { LogEntry } from '$lib/types';
import { useWebSocket } from '$lib/useWebSocket';

export const liveFilter = writable('');
export const paused = writable(false);
export const liveLogs = createBufferedLogsStore<LogEntry>({
	max: 1_000,
	flushInterval: 100
});
export const filteredLiveLogs = derived([liveLogs, liveFilter], ([$logs, $filter]) => {
	if (!$filter.trim()) return $logs;
	return $logs.filter((log) =>
		Object.values(log)
			.map((item) => {
				return typeof item === 'object' ? JSON.stringify(item) : String(item);
			})
			.join(' ')
			.toLowerCase()
			.includes($filter.toLowerCase())
	);
});

paused.subscribe((p) => liveLogs.setPaused(p));

if (browser) {
	useWebSocket(`ws://${location.host}/ws`, {
		onMessage: (data) => {
			const { timestamp, trace_id, level, message, ...rest } = data;
			const log: LogEntry = {
				timestamp,
				trace_id,
				level,
				message,
				raw: rest
			};
			liveLogs.add(log);
		}
	});
}
