import type { LayoutLoad } from './$types';

export const prerender = true;

export const load = (async () => {
	return {
		title: 'Magic Log UI'
	};
}) satisfies LayoutLoad;
