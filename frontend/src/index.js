import React from 'react';
import ReactDOM from 'react-dom/client';
import { QueryClient, QueryClientProvider } from 'react-query';
import App from './App';
import Websocket from './utils/ws';

const root = ReactDOM.createRoot(document.getElementById('root'));
const queryClient = new QueryClient();
Websocket.initialize(queryClient);

root.render(
	// This break react-virtualized for some reason
	// <React.StrictMode>
	<QueryClientProvider client={queryClient}>
		<App />
	</QueryClientProvider>
	// </React.StrictMode>
);
