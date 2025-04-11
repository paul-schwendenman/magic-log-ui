export function formatRelativeTime(timestamp: number): string {
	if (!Number.isFinite(timestamp)) return '';

	const now = Date.now();
	const diff = Math.floor((timestamp - now) / 1000); // in seconds

	const units: [Intl.RelativeTimeFormatUnit, number][] = [
		['year', 60 * 60 * 24 * 365],
		['month', 60 * 60 * 24 * 30],
		['day', 60 * 60 * 24],
		['hour', 60 * 60],
		['minute', 60],
		['second', 1]
	];

	const formatter = new Intl.RelativeTimeFormat(undefined, { numeric: 'auto' });

	for (const [unit, secondsInUnit] of units) {
		if (Math.abs(diff) >= secondsInUnit || unit === 'second') {
			return formatter.format(Math.round(diff / secondsInUnit), unit);
		}
	}

	return '';
}
