import { CssBaseline, StyledEngineProvider } from '@mui/material';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { RecoilRoot } from 'recoil';
import './App.css';
import Gallery from './components/Gallery';
import ItemPage from './components/ItemPage';
import TopBar from './components/TopBar';

function App() {
	return (
		<CssBaseline>
			<StyledEngineProvider injectFirst>
				<RecoilRoot>
					<BrowserRouter>
						<TopBar />
						<Routes>
							<Route index element={<Gallery />} />
							<Route path="/item/:itemId" element={<ItemPage />} />
						</Routes>
					</BrowserRouter>
				</RecoilRoot>
			</StyledEngineProvider>
		</CssBaseline>
	);
}

export default App;
