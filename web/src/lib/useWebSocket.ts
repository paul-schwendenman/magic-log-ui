import { writable } from 'svelte/store';

export const isConnected = writable(true);

export function useWebSocket(url: string, { onMessage }: { onMessage: (data: any) => void }) {
	let socket: WebSocket;
	let reconnectTimeout: ReturnType<typeof setTimeout>;

	function connect() {
		socket = new WebSocket(url);

		socket.onopen = () => {
			// console.log('WebSocket connected');
			isConnected.set(true);
		};

		socket.onmessage = (event) => {
			try {
				const data = JSON.parse(event.data);
				onMessage(data);
			} catch (err) {
				console.error('Error parsing WS message', err);
			}
		};

		socket.onclose = () => {
			// console.log('WebSocket disconnected. Reconnecting...');
			isConnected.set(false);
			reconnectTimeout = setTimeout(connect, 1000);
		};

		socket.onerror = (err) => {
			// console.error('WebSocket error', err);
			socket.close();
		};
	}

	connect();

	return {
		close: () => {
			clearTimeout(reconnectTimeout);
			socket?.close();
		}
	};
}
