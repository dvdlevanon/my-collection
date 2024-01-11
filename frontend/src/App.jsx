import { CssBaseline, StyledEngineProvider, ThemeProvider } from '@mui/material';
import React, { useState } from 'react';
import { useQuery } from 'react-query';
import { ReactQueryDevtools } from 'react-query/devtools';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import './App.css';
import Gallery from './components/pages/Gallery';
import ItemPage from './components/pages/ItemPage';
import ManageDirectories from './components/pages/ManageDirectories';
import TopBar from './components/top-bar/TopBar';
import Client from './utils/client';
import ReactQueryUtil from './utils/react-query-util';
import TagsUtil from './utils/tags-util';
import ThemeUtil from './utils/theme-utils';

function App() {
	useQuery({
		queryKey: ReactQueryUtil.SPECIAL_TAGS_KEY,
		queryFn: Client.getSpecialTags,
		onSuccess: (specialTags) => TagsUtil.initSpecialTags(specialTags),
	});

	useQuery({
		queryKey: ReactQueryUtil.CATEGORIES_KEY,
		queryFn: Client.getCategories,
		onSuccess: (categories) => TagsUtil.initCategories(categories),
	});

	const [hideTopBar, setHideTopBar] = useState(false);
	const [previewMode, setPreviewMode] = useState(true);

	return (
		<React.Fragment>
			<StyledEngineProvider injectFirst>
				<ThemeProvider theme={ThemeUtil.createDarkTheme()}>
					<CssBaseline enableColorScheme />
					<BrowserRouter>
						{!hideTopBar && <TopBar previewMode={previewMode} onPreviewModeChange={setPreviewMode} />}
						<Routes>
							<Route
								index
								element={<Gallery previewMode={previewMode} setHideTopBar={setHideTopBar} />}
							/>
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
