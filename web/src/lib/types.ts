export type LogEntry = {
	id?: string;
	created_at?: string;
	trace_id?: string;
	level?: string;
	message?: string;
	raw?: any;
};

export type TimeRange = {
	from: Date;
	to: Date;
	label: string; // e.g. "Past 15 Minutes"
	durationMs?: number; // for live shifting
	live?: boolean;
};
