import { writable } from 'svelte/store';

export const isConnected = writable(true);
export function useWebSocket(url: string, { onMessage }: { onMessage: (data: any) => void }) {
	let socket: WebSocket;
	let reconnectTimeout: ReturnType<typeof setTimeout>;
	let retryDelay = 1000;

	function connect() {
		socket = new WebSocket(url);

		socket.onopen = () => {
			console.log('WebSocket connected');
			retryDelay = 1000;
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
			console.log('WebSocket disconnected. Reconnecting...');
			isConnected.set(false);
			reconnect();
		};

		socket.onerror = (err) => {
			if (socket.readyState !== WebSocket.CLOSED && socket.readyState !== WebSocket.CLOSING) {
				console.warn('WebSocket encountered an error (will reconnect):', err);
			}
			socket.close();
		};
	}
	function reconnect() {
		reconnectTimeout = setTimeout(() => {
			retryDelay = Math.min(retryDelay * 2, 10000);
			connect();
		}, retryDelay);
	}

	connect();

	return {
		close: () => {
			clearTimeout(reconnectTimeout);
			socket?.close();
		}
	};
}
