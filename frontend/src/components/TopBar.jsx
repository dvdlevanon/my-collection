import { Button } from '@mui/material';
import { Link } from 'react-router-dom';
import Client from '../network/client';
import styles from './TopBar.module.css';

function TopBar() {
	return (
		<div className={styles.top_bar}>
			<Link to="/">
				<Button variant="outlined">Home</Button>
			</Link>
			<Button variant="outlined" onClick={() => Client.refreshCovers()}>
				Refresh Covers
			</Button>
			<Button variant="outlined" onClick={() => Client.refreshPreview()}>
				Refresh Preview
			</Button>
		</div>
	);
}

export default TopBar;
