<script lang="ts">
	import type { PageProps } from './$types';
	import LogTable from '$lib/components/LogTable.svelte';

	let { data }: PageProps = $props();
	const { logs } = data;
</script>

<h2 class="mb-2 text-xl font-bold">Log Table</h2>
<LogTable {logs} />
