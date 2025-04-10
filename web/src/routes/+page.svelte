<script lang="ts">
	import { onMount } from "svelte";

    let query = 'SELECT * FROM logs ORDER BY timestamp DESC LIMIT 100';
    let results: any[] = [];
    let wsLogs: any[] = [];
    let socket: WebSocket;
  
    const fetchQuery = async () => {
      const res = await fetch(`/query?q=${encodeURIComponent(query)}`);
      results = await res.json();
    };
  
    onMount(() => {
        socket = new WebSocket(`ws://${location.host}/ws`);
        socket.onmessage = (msg) => {
        wsLogs = [JSON.parse(msg.data), ...wsLogs].slice(0, 100);
        };
    });
  </script>
  
  <div class="space-y-4">
    <input bind:value={query} class="w-full p-2 bg-gray-800 border border-gray-600 rounded" />
    <button on:click={fetchQuery} class="px-4 py-2 bg-blue-600 rounded hover:bg-blue-500">Run Query</button>
  
    <h2 class="text-xl font-bold">Live Logs</h2>
    {#each wsLogs as log}
      <pre class="text-sm bg-gray-800 p-2 rounded mb-1">{JSON.stringify(log, null, 2)}</pre>
    {/each}
  
    <h2 class="text-xl font-bold mt-6">Query Results</h2>
    {#each results as row}
      <pre class="text-sm bg-gray-700 p-2 rounded mb-1">{JSON.stringify(row, null, 2)}</pre>
    {/each}
  </div>
  