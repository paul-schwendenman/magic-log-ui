import { derived, writable } from 'svelte/store';

export const wsStatus = writable<'connecting' | 'open' | 'closed' | 'error'>('connecting');
export const isConnected = derived(wsStatus, ($status) => $status === 'open');

export function useWebSocket(url: string, { onMessage }: { onMessage: (data: any) => void }) {
	let socket: WebSocket;
	let reconnectTimeout: ReturnType<typeof setTimeout>;
	let retryDelay = 1000;

	function connect() {
		socket = new WebSocket(url);

		socket.onopen = () => {
			console.log('WebSocket connected');
			retryDelay = 1000;
			wsStatus.set('open');
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
			wsStatus.set('closed');
			reconnect();
		};

		socket.onerror = (err) => {
			if (socket.readyState !== WebSocket.CLOSED && socket.readyState !== WebSocket.CLOSING) {
				console.warn('WebSocket encountered an error (will reconnect):', err);
			}
			wsStatus.set('error');
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
