<script lang="ts">
	import PresetEditor from '$lib/components/PresetEditor.svelte';
	import type { Config, ConfigDefaults, ConfigFieldTypes } from '$lib/types';
	import { typedEntries } from '$lib/utils';
	import { onMount } from 'svelte';

	let config: Config | null = null;
	let loading = true;
	let saving = false;

	let errorMessages: string[] = [];
	let successMessage = '';
	let successTimeout: ReturnType<typeof setTimeout> | null = null;

	const configFieldTypes: ConfigFieldTypes = {
		db_file: 'string',
		port: 'number',
		launch: 'boolean',
		log_format: 'string',
		regex_preset: 'preset',
		regex: 'string',
		jq: 'string',
		jq_preset: 'preset',
		csv_fields: 'string',
		has_csv_header: 'boolean'
	};

	const defaultConfig: ConfigDefaults = {
		jq_presets: {},
		regex_presets: {},
		has_csv_header: false,
		launch: false,
		port: 3000
	};

	$: regexPresetOptions = config?.regex_presets ? Object.keys(config.regex_presets) : [];
	$: jqPresetOptions = config?.jq_presets ? Object.keys(config.jq_presets) : [];

	onMount(async () => {
		try {
			const res = await fetch('/api/config');
			const rawConfig: Partial<Config> = await res.json();
			config = { ...defaultConfig, ...rawConfig };
		} catch (e) {
			errorMessages = ['Failed to load config.'];
		} finally {
			loading = false;
		}
	});

	async function save() {
		saving = true;
		errorMessages = [];
		successMessage = '';
		if (successTimeout) clearTimeout(successTimeout);

		try {
			const res = await fetch('/api/config', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(config)
			});

			if (!res.ok) {
				const contentType = res.headers.get('Content-Type');
				if (contentType?.includes('application/json')) {
					const data = await res.json();
					errorMessages = Array.isArray(data.errors) ? data.errors : [JSON.stringify(data)];
				} else {
					const text = await res.text();
					errorMessages = [text || 'Unknown error'];
				}
				return;
			}

			successMessage = '✅ Config saved successfully!';
			successTimeout = setTimeout(() => {
				successMessage = '';
			}, 4000);
		} catch (e) {
			errorMessages = ['Failed to save config.'];
		} finally {
			saving = false;
		}
	}

	function validateRegex(value: string): string | null {
		try {
			new RegExp(value);
			return null;
		} catch {
			return 'Invalid regular expression.';
		}
	}

	function validateJQ(value: string): string | null {
		if (!value.trim()) return 'JQ filter cannot be empty.';
		// For real JQ parsing/validation, do it on backend
		return null;
	}
</script>

{#if loading}
	<p class="p-4">Loading...</p>
{:else}
	<div class="mx-auto max-w-3xl space-y-6 p-4">
		<h1 class="text-2xl font-bold">App Settings</h1>

		{#if errorMessages.length}
			<ul class="space-y-1 rounded bg-red-100 p-3 text-red-700">
				{#each errorMessages as msg}
					<li>• {msg}</li>
				{/each}
			</ul>
		{/if}

		{#if successMessage}
			<div class="rounded bg-green-100 p-3 text-green-700">{successMessage}</div>
		{/if}

		{#if config}
			<section>
				<h2 class="mb-2 text-xl font-semibold">Defaults</h2>
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
					{#each typedEntries(configFieldTypes) as [key, type]}
						<div>
							<label for={key} class="block font-medium capitalize">{key}</label>
							{#if type === 'boolean'}
								<select id={key} class="w-full rounded border p-2" bind:value={config[key]}>
									<option value={true}>true</option>
									<option value={false}>false</option>
								</select>
							{:else if type === 'number'}
								<input
									id={key}
									type="number"
									class="w-full rounded border p-2"
									bind:value={config[key]}
									on:input={(e) => {
										const value = (e.target as HTMLInputElement).value;
										(config as any)[key] = +value;
									}}
								/>
							{:else if type === 'preset'}
								<select id={key} class="w-full rounded border p-2" bind:value={config[key]}>
									<option value="">-- select --</option>
									{#each key === 'regex_preset' ? regexPresetOptions : jqPresetOptions as option}
										<option value={option}>{option}</option>
									{/each}
								</select>
							{:else}
								<input id={key} class="w-full rounded border p-2" bind:value={config[key]} />
							{/if}
						</div>
					{/each}
				</div>
			</section>

			<PresetEditor
				title="Regex Presets"
				bind:presets={config.regex_presets}
				validateValue={validateRegex}
			/>

			<PresetEditor
				title="JQ Presets"
				bind:presets={config.jq_presets}
				validateValue={validateJQ}
			/>
		{/if}

		<button
			on:click={save}
			class="rounded bg-blue-600 px-4 py-2 text-white hover:bg-blue-700 disabled:opacity-50"
			disabled={saving}
		>
			{saving ? 'Saving...' : 'Save Changes'}
		</button>
	</div>
{/if}
