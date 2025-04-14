<script lang="ts">
	import Toast from '$lib/components/Toast.svelte';

	export let value: any;
	let expanded = false;
	let showToast = false;

	$: formatted = expanded ? JSON.stringify(value, null, 2) : JSON.stringify(value);

	function toggle() {
		expanded = !expanded;
	}

	async function copy() {
		await navigator.clipboard.writeText(formatted);
		showToast = true;
		setTimeout(() => (showToast = false), 1600);
	}
</script>

<div class="group relative cursor-pointer text-left whitespace-pre-wrap" on:click={toggle}>
	<pre class="max-w-full overflow-x-auto pr-6">{formatted}</pre>

	<button
		class="absolute top-1 right-1 hidden rounded bg-gray-700 p-1 text-white group-hover:block hover:bg-gray-600"
		on:click|stopPropagation={copy}
		aria-label="Copy raw log"
	>
		<svg
			xmlns="http://www.w3.org/2000/svg"
			fill="none"
			viewBox="0 0 24 24"
			stroke-width="1.5"
			stroke="currentColor"
			class="size-4"
		>
			<path
				stroke-linecap="round"
				stroke-linejoin="round"
				d="M15.666 3.888A2.25 2.25 0 0 0 13.5 2.25h-3c-1.03 0-1.9.693-2.166 1.638m7.332 0c.055.194.084.4.084.612v0a.75.75 0 0 1-.75.75H9a.75.75 0 0 1-.75-.75v0c0-.212.03-.418.084-.612m7.332 0c.646.049 1.288.11 1.927.184 1.1.128 1.907 1.077 1.907 2.185V19.5a2.25 2.25 0 0 1-2.25 2.25H6.75A2.25 2.25 0 0 1 4.5 19.5V6.257c0-1.108.806-2.057 1.907-2.185a48.208 48.208 0 0 1 1.927-.184"
			/>
		</svg>
	</button>

	{#if showToast}
		<div class="absolute -top-8 right-0">
			<Toast text="Copied to clipboard!" />
		</div>
	{/if}
</div>
