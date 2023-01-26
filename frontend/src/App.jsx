import { CssBaseline, StyledEngineProvider } from '@mui/material';
import React, { useState } from 'react';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { RecoilRoot } from 'recoil';
import './App.css';
import Gallery from './components/Gallery';
import ItemPage from './components/ItemPage';
import TopBar from './components/TopBar';

function App() {
	let [previewMode, setPreviewMode] = useState(true);
	const onPreviewModeChange = (previewMode) => {
		setPreviewMode(previewMode);
	};

	return (
		<React.Fragment>
			<CssBaseline />
			<StyledEngineProvider injectFirst>
				<RecoilRoot>
					<BrowserRouter>
						<TopBar previewMode={previewMode} onPreviewModeChange={onPreviewModeChange} />
						<Routes>
							<Route index element={<Gallery previewMode={previewMode} />} />
							<Route path="/spa/item/:itemId" element={<ItemPage />} />
						</Routes>
					</BrowserRouter>
				</RecoilRoot>
			</StyledEngineProvider>
		</React.Fragment>
	);
}

export default App;
