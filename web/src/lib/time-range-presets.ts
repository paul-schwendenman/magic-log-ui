import { sub } from 'date-fns';
import type { TimeRangeConfig } from './types';
import { m } from './paraglide/messages';

export const presets: TimeRangeConfig[] = [
	{
		label: m.grassy_wacky_okapi_clip(),
		from: sub(new Date(), { minutes: 5 }),
		to: new Date(),
		durationMs: 5 * 60_000,
		relative: true,
		live: true
	},
	{
		label: m.same_salty_marmot_grow(),
		from: sub(new Date(), { minutes: 15 }),
		to: new Date(),
		durationMs: 15 * 60_000,
		relative: true
	},
	{
		label: m.round_slow_goose_drop(),
		from: sub(new Date(), { minutes: 30 }),
		to: new Date(),
		durationMs: 30 * 60_000,
		relative: true
	},
	{
		label: m.flaky_major_fly_rise(),
		from: sub(new Date(), { hours: 1 }),
		to: new Date(),
		durationMs: 60 * 60_000,
		refreshMs: 30_000,
		relative: true
	},
	{
		label: m.patient_solid_slug_sail(),
		from: sub(new Date(), { hours: 4 }),
		to: new Date(),
		durationMs: 4 * 60 * 60_000,
		refreshMs: 30_000,
		relative: true
	},
	{
		label: m.pretty_spicy_okapi_scold(),
		from: sub(new Date(), { hours: 12 }),
		to: new Date(),
		durationMs: 12 * 60 * 60_000,
		refreshMs: 60_000,
		relative: true
	},
	{
		label: m.cute_flaky_pigeon_embrace(),
		from: sub(new Date(), { days: 1 }),
		to: new Date(),
		durationMs: 24 * 60 * 60_000,
		refreshMs: 60_000,
		relative: true
	},
	{
		label: m.honest_aqua_walrus_quell(),
		from: sub(new Date(), { days: 3 }),
		to: new Date(),
		durationMs: 3 * 24 * 60 * 60_000,
		refreshMs: 60_000,
		relative: true
	}
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
