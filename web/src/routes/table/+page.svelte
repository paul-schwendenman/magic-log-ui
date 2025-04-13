<script lang="ts">
	import { createColumnHelper, createSvelteTable, getCoreRowModel } from '$lib/table';
	// import * as UserService from '$lib/services/user-profile';
	import type { PageData } from './$types';

	type UserProfile = {
		id: string;
		name: string;
		age: number;
		email: string;
		phone: string;
		birthdate: string;
		friends: UserProfile[];
	};

	let { data }: PageProps<PageData> = $props();
	const {
		userProfiles = [
			{
				name: 'A',
				age: 'B',
				email: 'C',
				phone: 'D'
			}
		]
	} = data;

	// Create a column helper for the user profile data.
	// It's not necessary, but it helps with type stuff.
	const colHelp = createColumnHelper<UserProfile>();

	// Define the columns using the column helper.
	// This is a basic example. Check other examples for more complexity.
	const columnDefs = [
		colHelp.accessor('name', { header: 'Name' }),
		colHelp.accessor('age', { header: 'Age' }),
		colHelp.accessor('email', { header: 'Email' }),
		colHelp.accessor('phone', { header: 'Phone' })
	];

	// Create the table.
	const table = createSvelteTable({
		data: userProfiles,
		columns: columnDefs,
		getCoreRowModel: getCoreRowModel()
	});
</script>

<h2>Table Demo</h2>

<hr />

<table>
	<thead>
		<tr>
			{#each table.getHeaderGroups() as headerGroup}
				{#each headerGroup.headers as header}
					<th>{header.column.columnDef.header}</th>
				{/each}
			{/each}
		</tr>
	</thead>
	<tbody>
		{#each table.getRowModel().rows as row}
			<tr>
				{#each row.getVisibleCells() as cell}
					<td>{cell.getValue()}</td>
				{/each}
			</tr>
		{/each}
	</tbody>
</table>
