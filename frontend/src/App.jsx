import './App.css';
import { BrowserRouter, Route, Link, Routes } from "react-router-dom";
import { RecoilRoot } from "recoil";
import Gallery from './components/Gallery';
import ItemPage from './components/ItemPage';
import TopBar from './components/TopBar';
import { CssBaseline, StyledEngineProvider } from '@mui/material';

function App() {
	return (
		<CssBaseline>
			<StyledEngineProvider injectFirst>
				<TopBar />
				<RecoilRoot>
					<BrowserRouter>
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
