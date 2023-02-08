import { AppBar, Button, Checkbox, FormControlLabel, Link, Toolbar, Typography } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';
import Client from '../../network/client';

function TopBar({ previewMode, onPreviewModeChange }) {
	const previewsChange = (e) => {
		onPreviewModeChange(e.target.checked);
	};

	return (
		<AppBar position="static">
			<Toolbar>
				<Link sx={{ flexGrow: 1 }} component={RouterLink} to="/">
					<Typography variant="h5">My Collection</Typography>
				</Link>
				<Link component={RouterLink} to="spa/manage-directories">
					<Button variant="outlined">Manage Directories</Button>
				</Link>
				<Link href={Client.getExportMetadataUrl()} download>
					<Button variant="outlined">Export metadata</Button>
				</Link>
				<Button variant="outlined" onClick={() => Client.refreshCovers()}>
					Refresh Covers
				</Button>
				<Button variant="outlined" onClick={() => Client.refreshPreview()}>
					Refresh Preview
				</Button>
				<Button variant="outlined" onClick={() => Client.refreshVideoMetadata()}>
					Refresh Video Metadata
				</Button>
				<FormControlLabel
					label="Use Previews"
					control={<Checkbox checked={previewMode} onChange={(e) => previewsChange(e)} />}
				/>
			</Toolbar>
		</AppBar>
	);
}

export default TopBar;
