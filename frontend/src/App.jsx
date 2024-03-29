import { CssBaseline, Stack, StyledEngineProvider, ThemeProvider } from '@mui/material';
import React, { useEffect, useState } from 'react';
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
	const [theme, setTheme] = useState(ThemeUtil.createDarkOrangeTheme());

	useEffect(() => {
		let themeName = localStorage.getItem('theme');
		if (themeName) {
			let theme = ThemeUtil.themeByName(themeName);

			if (theme) {
				setTheme(ThemeUtil.themeByName(themeName));
			}
		}
	}, []);

	return (
		<React.Fragment>
			<StyledEngineProvider injectFirst>
				<ThemeProvider theme={theme}>
					<CssBaseline enableColorScheme />
					<Stack
						sx={{
							backgroundImage: theme.backgroundImage,
							width: '100%',
							height: '100%',
						}}
					>
						<BrowserRouter>
							{!hideTopBar && (
								<TopBar
									previewMode={previewMode}
									onPreviewModeChange={setPreviewMode}
									theme={theme}
									setTheme={(theme) => {
										if (!theme) {
											return;
										}

										setTheme(theme);
										localStorage.setItem('theme', theme.name);
									}}
								/>
							)}
							<Routes>
								<Route
									index
									element={<Gallery previewMode={previewMode} setHideTopBar={setHideTopBar} />}
								/>
								<Route path="/spa/item/:itemId" element={<ItemPage />} />
								<Route path="/spa/manage-directories" element={<ManageDirectories />} />
							</Routes>
						</BrowserRouter>
					</Stack>
					<ReactQueryDevtools initialIsOpen={false} />
				</ThemeProvider>
			</StyledEngineProvider>
		</React.Fragment>
	);
}

export default App;
