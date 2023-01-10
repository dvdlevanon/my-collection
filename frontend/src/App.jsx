import './App.css';
import { BrowserRouter, Route, Link, Routes } from "react-router-dom";
import Gallery from './pages/Gallery';
import ItemPage from './pages/ItemPage';
import TopBar from './components/TopBar';

function App() {
	return (
		<>
			<TopBar />
			<BrowserRouter>
				<Routes>
					<Route index element={<Gallery />} />
					<Route path="/item/:itemId" element={<ItemPage />} />
				</Routes>
			</BrowserRouter>
		</>
	);
}

export default App;
