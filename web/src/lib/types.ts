export type LogEntry = {
	id?: string;
	timestamp?: string;
	trace_id?: string;
	level?: string;
	message?: string;
	raw?: any;
};

export type TimeRange = {
	from: Date;
	to: Date;
};

export type TimeRangeConfig = {
	from: Date;
	to: Date;
	label: string; // e.g. "Past 15 Minutes"
	durationMs?: number; // for live shifting
	refreshMs?: number;
	live?: boolean;
	relative?: boolean;
};

export type Config = {
	db_file?: string;
	port: number;
	launch: boolean;
	log_format?: string;
	regex_preset?: string;
	regex?: string;
	jq?: string;
	jq_preset?: string;
	csv_fields?: string;
	has_csv_header: boolean;
	regex_presets: Record<string, string>;
	jq_presets: Record<string, string>;
};

export type ConfigDefaults = Pick<
	Config,
	'regex_presets' | 'jq_presets' | 'has_csv_header' | 'launch' | 'port'
>;

type ConfigKey = Exclude<keyof Config, 'jq_presets' | 'regex_presets'>;

export type ConfigValueType = 'string' | 'number' | 'boolean' | 'preset';

export type ConfigFieldTypes = Record<ConfigKey, ConfigValueType>;
