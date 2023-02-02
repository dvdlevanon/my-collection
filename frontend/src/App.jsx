import { createTheme, CssBaseline, StyledEngineProvider, ThemeProvider } from '@mui/material';
import React, { useState } from 'react';
import { QueryClient, QueryClientProvider } from 'react-query';
import { ReactQueryDevtools } from 'react-query/devtools';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import './App.css';
import Gallery from './components/Gallery';
import ItemPage from './components/ItemPage';
import ManageDirectories from './components/ManageDirectories';
import TopBar from './components/TopBar';
const theme = createTheme({
	palette: {
		mode: 'dark',
		primary: {
			main: '#ff4400',
		},
		secondary: {
			main: '#0D9352',
		},
		bright: {
			main: '#ffddcc',
		},
		dark: {
			main: '#554433',
		},
	},
	typography: {
		fontSize: 16,
	},
});

function App() {
	const queryClient = new QueryClient();
	let [previewMode, setPreviewMode] = useState(true);
	const onPreviewModeChange = (previewMode) => {
		setPreviewMode(previewMode);
	};

	return (
		<React.Fragment>
			<StyledEngineProvider injectFirst>
				<ThemeProvider theme={theme}>
					<QueryClientProvider client={queryClient}>
						<CssBaseline />
						<BrowserRouter>
							<TopBar previewMode={previewMode} onPreviewModeChange={onPreviewModeChange} />
							<Routes>
								<Route index element={<Gallery previewMode={previewMode} />} />
								<Route path="/spa/item/:itemId" element={<ItemPage />} />
								<Route path="/spa/manage-directories" element={<ManageDirectories />} />
							</Routes>
						</BrowserRouter>
						<ReactQueryDevtools initialIsOpen={false} />
					</QueryClientProvider>
				</ThemeProvider>
			</StyledEngineProvider>
		</React.Fragment>
	);
}

export default App;
