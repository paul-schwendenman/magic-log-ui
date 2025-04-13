<script lang="ts">
	import { createColumnHelper, createSvelteTable, getCoreRowModel } from '$lib/table';
	import type { PageProps } from '../$types';
	import type { PageData } from './$types';

	type LogEntry = {
		timestamp?: string;
		trace_id?: string;
		level?: string;
		message?: string;
		raw?: any;
	};

	let { data }: PageProps<PageData> = $props();
	const { logs } = data;

	const colHelp = createColumnHelper<LogEntry>();

	const columnDefs = [
		colHelp.accessor('timestamp', { header: 'Time' }),
		colHelp.accessor('trace_id', { header: 'Trace ID' }),
		colHelp.accessor('level', { header: 'Level' }),
		colHelp.accessor('message', { header: 'Message' })
	];

	const table = createSvelteTable({
		data: logs,
		columns: columnDefs,
		getCoreRowModel: getCoreRowModel()
	});
</script>

<h2 class="text-xl font-bold mb-2">Log Table</h2>
<table class="min-w-full text-sm border-collapse border">
	<thead>
		<tr>
			{#each table.getHeaderGroups() as headerGroup}
				{#each headerGroup.headers as header}
					<th class="border px-2 py-1 font-semibold">{header.column.columnDef.header}</th>
				{/each}
			{/each}
		</tr>
	</thead>
	<tbody>
		{#each table.getRowModel().rows as row}
			<tr>
				{#each row.getVisibleCells() as cell}
					<td class="border px-2 py-1">{cell.getValue()}</td>
				{/each}
			</tr>
		{/each}
	</tbody>
</table>
