import type { PageLoad } from './$types';

export const load = (async ({ fetch }) => {
	const logs = [];
	// try {
	// 	const res = await fetch('/query?q=SELECT * FROM logs ORDER BY timestamp DESC LIMIT 100');
	// 	if (res.ok) {
	// 		logs = await res.json();
	// 	}
	// } catch {}

	return {
		logs
	};
}) satisfies PageLoad;
