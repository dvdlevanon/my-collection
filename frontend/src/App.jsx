import { createTheme, CssBaseline, StyledEngineProvider, ThemeProvider } from '@mui/material';
import React, { useState } from 'react';
import { ReactQueryDevtools } from 'react-query/devtools';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import './App.css';
import Gallery from './components/pages/Gallery';
import ItemPage from './components/pages/ItemPage';
import ManageDirectories from './components/pages/ManageDirectories';
import TopBar from './components/top-bar/TopBar';

const theme = createTheme({
	palette: {
		mode: 'dark',
		primary: {
			main: '#ff4400',
			light: '#ff8844',
		},
		secondary: {
			main: '#0D9352',
		},
		bright: {
			main: '#ffddcc',
			darker: '#ddbbaa',
		},
		dark: {
			main: '#121212',
			lighter: '#222222',
			lighter2: '#272727',
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
							<Route path="/spa/manage-directories" element={<ManageDirectories />} />
						</Routes>
					</BrowserRouter>
					<ReactQueryDevtools initialIsOpen={false} />
				</ThemeProvider>
			</StyledEngineProvider>
		</React.Fragment>
	);
}

export default App;
