<script lang="ts">
	import { m } from '$lib/paraglide/messages.js';
	import JsonInline from './JsonInline.svelte';
	import JsonViewer from './JsonViewer.svelte';
	export let log: any;

	let expanded = false;
	let copied = false;

	function toggle() {
		expanded = !expanded;
	}

	function expand() {
		if (!expanded) {
			expanded = true;
		}
	}

	async function copy() {
		const text = JSON.stringify(log);

		try {
			await navigator.clipboard.writeText(text);
			copied = true;
			setTimeout(() => (copied = false), 1000);
		} catch (err) {
			console.error('Failed to copy:', err);
		}
	}
</script>

<div
	class="flex items-center gap-1 rounded bg-gray-800 p-2 font-mono text-sm whitespace-pre-wrap transition-all hover:bg-gray-700"
>
	<div class="order-last flex self-start">
		{#if copied}
			<span class="relative text-xs text-green-400">{m.copied()}</span>
		{/if}
		<button onclick={copy} aria-label="copy">
			<svg
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
				stroke-width="1.5"
				stroke="currentColor"
				class="size-6"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="M15.666 3.888A2.25 2.25 0 0 0 13.5 2.25h-3c-1.03 0-1.9.693-2.166 1.638m7.332 0c.055.194.084.4.084.612v0a.75.75 0 0 1-.75.75H9a.75.75 0 0 1-.75-.75v0c0-.212.03-.418.084-.612m7.332 0c.646.049 1.288.11 1.927.184 1.1.128 1.907 1.077 1.907 2.185V19.5a2.25 2.25 0 0 1-2.25 2.25H6.75A2.25 2.25 0 0 1 4.5 19.5V6.257c0-1.108.806-2.057 1.907-2.185a48.208 48.208 0 0 1 1.927-.184"
				/>
			</svg>
		</button>
		<button onclick={toggle} aria-label="toggle">
			{#if !expanded}
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="1.5"
					stroke="currentColor"
					class="size-6"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M3.75 3.75v4.5m0-4.5h4.5m-4.5 0L9 9M3.75 20.25v-4.5m0 4.5h4.5m-4.5 0L9 15M20.25 3.75h-4.5m4.5 0v4.5m0-4.5L15 9m5.25 11.25h-4.5m4.5 0v-4.5m0 4.5L15 15"
					/>
				</svg>
			{:else}
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="1.5"
					stroke="currentColor"
					class="size-6"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M9 9V4.5M9 9H4.5M9 9 3.75 3.75M9 15v4.5M9 15H4.5M9 15l-5.25 5.25M15 9h4.5M15 9V4.5M15 9l5.25-5.25M15 15h4.5M15 15v4.5m0-4.5 5.25 5.25"
					/>
				</svg>
			{/if}
		</button>
	</div>
	<div class="flex-grow truncate" onclick={expand} aria-label="expand">
		{#if expanded}
			<JsonViewer data={log} />
		{:else}
			<JsonInline data={log} />
		{/if}
	</div>
</div>
