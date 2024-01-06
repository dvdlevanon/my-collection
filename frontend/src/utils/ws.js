import Client from './client';
import ReactQueryUtil from './react-query-util';

export default class Websocket {
	static socket = null;
	static queryClient = null;

	// values in types.go
	static PUSH_PING = 1;
	static PUSH_QUEUE_METADATA = 2;

	static initialize = (queryClient) => {
		Websocket.socket = new WebSocket(Client.websocketUrl);
		Websocket.queryClient = queryClient;

		Websocket.socket.addEventListener('open', (event) => {});

		Websocket.socket.addEventListener('message', (event) => {
			let message = JSON.parse(event.data);

			if (message.type == Websocket.PUSH_QUEUE_METADATA) {
				Websocket.queryClient.setQueryData(ReactQueryUtil.QUEUE_METADATA_KEY, message.payload);
			} else if (message.type == Websocket.PUSH_PING) {
			}
		});
	};
}
