<script lang="ts">
	import ExpandableCell from '$lib/components/ExpandableCell.svelte';
	import {
		createColumnHelper,
		createSvelteTable,
		FlexRender,
		getCoreRowModel,
		renderComponent
	} from '$lib/table';
	import type { LogEntry } from '$lib/types';

	type Props = {
		logs: LogEntry[];
		initialVisibility?: Record<string, boolean>;
	};

	let { logs, initialVisibility = {} }: Props = $props();
	let dataState = $state(logs);
	let columnVisibility = $state({ ...initialVisibility });

	const colHelp = createColumnHelper<LogEntry>();

	const columnDefs = [
		colHelp.accessor('timestamp', { header: 'Time' }),
		colHelp.accessor('trace_id', { header: 'Trace ID' }),
		colHelp.accessor('level', { header: 'Level' }),
		colHelp.accessor('message', { header: 'Message' }),
		colHelp.accessor('raw', {
			header: 'Raw',
			cell: ({ cell }) => renderComponent(ExpandableCell, { value: cell.getValue() })
		})
	];

	const table = createSvelteTable({
		get data() {
			return logs;
		},
		columns: columnDefs,
		get state() {
			return {
				columnVisibility
			};
		},
		getCoreRowModel: getCoreRowModel()
	});
</script>

<table class="min-w-full border-collapse border text-sm">
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
					<td class="border px-2 py-1">
						<FlexRender content={cell.column.columnDef.cell} context={cell.getContext()} />
					</td>
				{/each}
			</tr>
		{/each}
	</tbody>
</table>
