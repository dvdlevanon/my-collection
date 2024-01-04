import { createTheme, CssBaseline, StyledEngineProvider, ThemeProvider } from '@mui/material';
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
			darker2: '#aaaaaa',
			text: '#ffffff',
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
	const specialTagsQuery = useQuery({
		queryKey: ReactQueryUtil.SPECIAL_TAGS_KEY,
		queryFn: Client.getSpecialTags,
		onSuccess: (specialTags) => TagsUtil.initSpecialTags(specialTags),
	});

	const categoriesQuery = useQuery({
		queryKey: ReactQueryUtil.CATEGORIES_KEY,
		queryFn: Client.getCategories,
		onSuccess: (categories) => TagsUtil.initCategories(categories),
	});

	const [hideTopBar, setHideTopBar] = useState(false);
	const [previewMode, setPreviewMode] = useState(true);

	return (
		<React.Fragment>
			<StyledEngineProvider injectFirst>
				<ThemeProvider theme={theme}>
					<CssBaseline />
					<BrowserRouter>
						{!hideTopBar && <TopBar previewMode={previewMode} onPreviewModeChange={setPreviewMode} />}
						{specialTagsQuery.isSuccess && (
							<Routes>
								<Route
									index
									element={<Gallery previewMode={previewMode} setHideTopBar={setHideTopBar} />}
								/>
								<Route path="/spa/item/:itemId" element={<ItemPage />} />
								<Route path="/spa/manage-directories" element={<ManageDirectories />} />
							</Routes>
						)}
					</BrowserRouter>
					<ReactQueryDevtools initialIsOpen={false} />
				</ThemeProvider>
			</StyledEngineProvider>
		</React.Fragment>
	);
}

export default App;
