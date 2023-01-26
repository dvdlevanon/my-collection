import { createTheme, CssBaseline, StyledEngineProvider, ThemeProvider } from '@mui/material';
import React, { useState } from 'react';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import './App.css';
import Gallery from './components/Gallery';
import ItemPage from './components/ItemPage';
import TopBar from './components/TopBar';

const theme = createTheme({
	palette: {
		mode: 'dark',
		primary: {
			main: '#ff4400',
		},
	},
	typography: {
		fontSize: 16,
	},
});

function App() {
	let [previewMode, setPreviewMode] = useState(true);
	const onPreviewModeChange = (previewMode) => {
		setPreviewMode(previewMode);
	};

	return (
		<React.Fragment>
			<StyledEngineProvider injectFirst>
				<ThemeProvider theme={theme}>
					<CssBaseline />
					<BrowserRouter>
						<TopBar previewMode={previewMode} onPreviewModeChange={onPreviewModeChange} />
						<Routes>
							<Route index element={<Gallery previewMode={previewMode} />} />
							<Route path="/spa/item/:itemId" element={<ItemPage />} />
						</Routes>
					</BrowserRouter>
				</ThemeProvider>
			</StyledEngineProvider>
		</React.Fragment>
	);
}

export default App;
