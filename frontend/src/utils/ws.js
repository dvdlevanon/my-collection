import Client from './client';
import ReactQueryUtil from './react-query-util';

export default class Websocket {
	static socket = null;
	static queryClient = null;

	static initialize = (queryClient) => {
		Websocket.socket = new WebSocket(Client.websocketUrl);
		Websocket.queryClient = queryClient;

		Websocket.socket.addEventListener('open', (event) => {});

		Websocket.socket.addEventListener('message', (event) => {
			let message = JSON.parse(event.data);

			if (message.type == 1) {
				Websocket.queryClient.setQueryData(ReactQueryUtil.QUEUE_METADATA_KEY, message.payload);
			} else if (message == 0) {
				//ping
			}
		});
	};
}
