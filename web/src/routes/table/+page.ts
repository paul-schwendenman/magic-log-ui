import type { PageLoad } from './$types';

export const load = (async ({ fetch }) => {
	const res = await fetch('/query?q=SELECT * FROM logs ORDER BY timestamp DESC LIMIT 100');
	const logs = await res.json();
	return {
		logs
	};
}) satisfies PageLoad;
