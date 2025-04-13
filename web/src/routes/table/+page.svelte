<script lang="ts">
	import { onMount } from 'svelte';
	import {
		createTable,
		getCoreRowModel,
		flexRender,
		type ColumnDef,
		type Table
	} from '@tanstack/svelte-table';
	import { npm_lifecycle_event } from '$env/static/private';

	let data: any[] = [];
	let columns: ColumnDef<any>[] = [];
	let table: Table<any> | null = null;

	onMount(async () => {
		const res = await fetch('/query?q=SELECT * FROM logs ORDER BY timestamp DESC LIMIT 100');
		const result = await res.json();

		if (result.length > 0) {
			const keys = Object.keys(result[0]);
			columns = keys.map((key) => ({
				accessorKey: key,
				header: key
			}));
		}

		data = result;

		table = createTable({
			data,
			columns,
			getCoreRowModel: getCoreRowModel()
		});
	});
</script>

<table class="min-w-full text-sm">
	<thead>
		{#each table.getHeaderGroups() as headerGroup}
			<tr>
				{#each headerGroup.headers as header}
					<th class="border-b px-2 py-1 text-left font-medium">
						{#if !header.isPlaceholder}
							{flexRender(header.column.columnDef.header, header.getContext())}
						{/if}
					</th>
				{/each}
			</tr>
		{/each}
	</thead>
	<tbody>
		{#each table.getRowModel().rows as row}
			<tr>
				{#each row.getVisibleCells() as cell}
					<td class="border-b px-2 py-1"
						>{flexRender(cell.column.columnDef.cell, cell.getContext())}</td
					>
				{/each}
			</tr>
		{/each}
	</tbody>
</table>
