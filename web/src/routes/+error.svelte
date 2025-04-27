<script>
	import { page } from '$app/state';
</script>

<svelte:head>
	<title>{page.status} - {page.error?.message || 'Error'}</title>
</svelte:head>

<main class="flex min-h-screen flex-col items-center justify-center">
	<h1 class="mb-4 text-5xl font-bold">{page.status}</h1>
	<p class="mb-6 text-xl">{page.error?.message || 'Something went wrong'}</p>
	<a href="/" class="text-blue-500 hover:underline">Go back home</a>
</main>
