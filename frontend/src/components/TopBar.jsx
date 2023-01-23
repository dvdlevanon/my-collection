import { Button, Checkbox, FormControlLabel } from '@mui/material';
import { Link } from 'react-router-dom';
import Client from '../network/client';
import styles from './TopBar.module.css';

function TopBar({ previewMode, onPreviewModeChange }) {
	const previewsChange = (e) => {
		onPreviewModeChange(e.target.checked);
	};

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
			<FormControlLabel
				label="Use Previews"
				control={<Checkbox checked={previewMode} onChange={(e) => previewsChange(e)} />}
			/>
		</div>
	);
}

export default TopBar;
