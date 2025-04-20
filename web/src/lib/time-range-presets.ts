import { sub } from 'date-fns';
import type { TimeRange } from './types';

export const presets: TimeRange[] = [
	{
		label: 'Past 15 Minutes',
		from: sub(new Date(), { minutes: 15 }),
		to: new Date(),
		durationMs: 15 * 60_000
	},
	{
		label: 'Past 1 Hour',
		from: sub(new Date(), { hours: 1 }),
		to: new Date(),
		durationMs: 60 * 60_000
	},
	{
		label: 'Past 1 Day',
		from: sub(new Date(), { days: 1 }),
		to: new Date(),
		durationMs: 24 * 60 * 60_000
	}
	// ...
];
