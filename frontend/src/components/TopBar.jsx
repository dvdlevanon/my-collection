import { Button } from '@mui/material';
import { Link } from 'react-router-dom';
import Client from '../network/client';
import styles from './TopBar.module.css';

function TopBar() {
	const refreshClicked = (e) => {
		Client.refreshPreview();
	};

	return (
		<div className={styles.top_bar}>
			<Link to="/">
				<Button variant="outlined">Home</Button>
			</Link>
			<Button variant="outlined" onClick={(e) => refreshClicked(e)}>
				Refresh Gallery
			</Button>
		</div>
	);
}

export default TopBar;
