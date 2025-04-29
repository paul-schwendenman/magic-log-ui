<script lang="ts">
	import ExpandableCell from '$lib/components/ExpandableCell.svelte';
	import { m } from '$lib/paraglide/messages';
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
		colHelp.accessor('timestamp', {
			id: 'timestamp',
			header: m.born_each_lobster_sing(),
			cell: ({ cell }) => {
				const isoString = cell.getValue();
				const date = new Date(isoString);
				// Format nicely, e.g. "Apr 26, 10:15 PM"
				const friendly = date.toLocaleString(undefined, {
					month: 'short',
					day: 'numeric',
					hour: 'numeric',
					minute: '2-digit',
					second: '2-digit',
					hour12: true
				});
				return friendly;
			}
		}),

		colHelp.accessor('trace_id', { id: 'trace_id', header: m.tense_blue_maggot_nurture() }),
		colHelp.accessor('level', { id: 'level', header: m.bland_ok_dolphin_kiss() }),
		colHelp.accessor('message', { id: 'message', header: m.hour_top_cat_fear() }),
		colHelp.accessor('raw', {
			id: 'raw',
			header: m.agent_trite_cow_list(),
			cell: ({ cell }) => renderComponent(ExpandableCell, { value: cell.getValue() })
		})
	];

	const table = createSvelteTable({
		get data() {
			return logs;
		},
		columns: columnDefs,
		onColumnVisibilityChange: (updater) => {
			const next = typeof updater === 'function' ? updater(columnVisibility) : updater;
			columnVisibility = next;
		},
		get state() {
			return {
				columnVisibility
			};
		},
		getCoreRowModel: getCoreRowModel()
	});
</script>

<div class="mb-2 flex flex-wrap gap-2 text-sm">
	{#each table.getAllLeafColumns() as col (col.id)}
		<label class="flex items-center gap-1">
			<input
				type="checkbox"
				checked={col.getIsVisible()}
				disabled={!col.getCanHide()}
				onchange={() => col.toggleVisibility()}
			/>
			{col.columnDef.header}
		</label>
	{/each}
</div>

<table class="min-w-full border-collapse border text-sm">
	<thead>
		<tr>
			{#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
				{#each headerGroup.headers as header (header.column.id)}
					{#if columnVisibility[header.column.id]}
						<th class="border px-2 py-1 font-semibold">{header.column.columnDef.header}</th>
					{/if}
				{/each}
			{/each}
		</tr>
	</thead>
	<tbody>
		{#each table.getRowModel().rows as row}
			<tr>
				{#each row.getVisibleCells() as cell (cell.id)}
					<td class="border px-2 py-1">
						<FlexRender content={cell.column.columnDef.cell} context={cell.getContext()} />
					</td>
				{/each}
			</tr>
		{/each}
	</tbody>
</table>
