import { sub } from 'date-fns';
import type { TimeRange } from './types';

export const presets: TimeRange[] = [
	{
		label: '5 Minutes',
		from: sub(new Date(), { minutes: 5 }),
		to: new Date(),
		durationMs: 5 * 60_000,
		live: true
	},
	{
		label: '15 Minutes',
		from: sub(new Date(), { minutes: 15 }),
		to: new Date(),
		durationMs: 15 * 60_000,
		live: true
	},
	{
		label: '30 Minutes',
		from: sub(new Date(), { minutes: 30 }),
		to: new Date(),
		durationMs: 30 * 60_000,
		live: true
	},
	{
		label: '1 Hour',
		from: sub(new Date(), { hours: 1 }),
		to: new Date(),
		durationMs: 60 * 60_000,
		refreshMs: 30_000,
		live: true
	},
	{
		label: '4 Hours',
		from: sub(new Date(), { hours: 4 }),
		to: new Date(),
		durationMs: 4 * 60 * 60_000,
		refreshMs: 30_000,
		live: true
	},
	{
		label: '12 Hours',
		from: sub(new Date(), { hours: 12 }),
		to: new Date(),
		durationMs: 12 * 60 * 60_000,
		refreshMs: 60_000,
		live: true
	},
	{
		label: '24 Hours',
		from: sub(new Date(), { days: 1 }),
		to: new Date(),
		durationMs: 24 * 60 * 60_000,
		refreshMs: 60_000,
		live: true
	},
	{
		label: '3 Days',
		from: sub(new Date(), { days: 3 }),
		to: new Date(),
		durationMs: 3 * 24 * 60 * 60_000,
		refreshMs: 60_000,
		live: true
	},
	// {
	// 	label: '7 Days',
	// 	from: sub(new Date(), { days: 7 }),
	// 	to: new Date(),
	// 	durationMs: 7 * 24 * 60 * 60_000,
	// 	refreshMs: 60_000,
	// 	live: true
	// },
	// {
	// 	label: 'All Time',
	// 	from: null,
	// 	to: null,
	// 	durationMs: undefined,
	// 	live: false // this one shouldn’t refresh
	// }
];
