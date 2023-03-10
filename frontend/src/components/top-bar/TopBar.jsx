import QueueIcon from '@mui/icons-material/List';
import {
	AppBar,
	Badge,
	Button,
	Checkbox,
	FormControlLabel,
	IconButton,
	Link,
	Menu,
	MenuItem,
	Popover,
	Toolbar,
	Tooltip,
	Typography,
} from '@mui/material';
import { useState } from 'react';
import { useQuery } from 'react-query';
import { Link as RouterLink } from 'react-router-dom';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import Queue from '../queue/Queue';

function TopBar({ previewMode, onPreviewModeChange }) {
	const [refreshAnchorEl, setRefreshAnchorEl] = useState(null);
	const queueMetadataQuery = useQuery({
		queryKey: ReactQueryUtil.QUEUE_METADATA_KEY,
		queryFn: Client.getQueueMetadata,
		onSuccess: (queueMetadata) => {
			if (queueMetadata.size == 0) {
				setQueueEl(null);
			}
		},
	});

	const previewsChange = (e) => {
		onPreviewModeChange(e.target.checked);
	};

	const [queueEl, setQueueEl] = useState(null);

	return (
		<AppBar position="static">
			<Toolbar
				sx={{
					gap: '20px',
					alignItems: 'center',
					alignContent: 'center',
					verticalAlign: 'center',
				}}
			>
				<Link sx={{ flexGrow: 1 }} component={RouterLink} to="/">
					<Typography variant="h5">My Collection</Typography>
				</Link>
				{queueMetadataQuery.isSuccess && (
					<Tooltip
						title={
							queueMetadataQuery.data.size == 0
								? 'No tasks'
								: queueMetadataQuery.data.unfinishedTasks + ' pending tasks'
						}
					>
						<span>
							<IconButton
								disabled={queueMetadataQuery.data.size == 0}
								onClick={(e) => setQueueEl(e.currentTarget)}
							>
								<Badge badgeContent={queueMetadataQuery.data.unfinishedTasks || null} color="primary">
									<QueueIcon />
								</Badge>
							</IconButton>
						</span>
					</Tooltip>
				)}
				{queueMetadataQuery.isSuccess && queueMetadataQuery.data.size != 0 && Boolean(queueEl) && (
					<Popover
						id="test"
						anchorEl={queueEl}
						open={Boolean(queueEl)}
						onClose={() => setQueueEl(null)}
						anchorOrigin={{
							vertical: 'bottom',
							horizontal: 'left',
						}}
					>
						<Queue onClose={() => setQueueEl(null)} />
					</Popover>
				)}
				<Link component={RouterLink} to="spa/manage-directories">
					<Button variant="outlined">Manage Directories</Button>
				</Link>
				<Link href={Client.getExportMetadataUrl()} download>
					<Button variant="outlined">Export metadata</Button>
				</Link>
				<Button variant="outlined" onClick={(e) => setRefreshAnchorEl(e.currentTarget)}>
					Refresh
				</Button>
				<Menu
					anchorEl={refreshAnchorEl}
					open={refreshAnchorEl != null}
					onClose={() => setRefreshAnchorEl(null)}
				>
					<MenuItem
						onClick={() => {
							Client.refreshCovers(false);
							setRefreshAnchorEl(null);
						}}
					>
						Refresh Covers
					</MenuItem>
					<MenuItem
						onClick={() => {
							Client.refreshCovers(true);
							setRefreshAnchorEl(null);
						}}
					>
						Force Refresh Covers
					</MenuItem>
					<MenuItem
						onClick={() => {
							Client.refreshPreview(false);
							setRefreshAnchorEl(null);
						}}
					>
						Refresh Preview
					</MenuItem>
					<MenuItem
						onClick={() => {
							Client.refreshPreview(true);
							setRefreshAnchorEl(null);
						}}
					>
						Force Refresh Preview
					</MenuItem>
					<MenuItem
						onClick={() => {
							Client.refreshVideoMetadata(false);
							setRefreshAnchorEl(null);
						}}
					>
						Refresh Video Metadata
					</MenuItem>
					<MenuItem
						onClick={() => {
							Client.refreshVideoMetadata(true);
							setRefreshAnchorEl(null);
						}}
					>
						Force Refresh Video Metadata
					</MenuItem>
				</Menu>
				<FormControlLabel
					label="Use Previews"
					control={<Checkbox checked={previewMode} onChange={(e) => previewsChange(e)} />}
				/>
			</Toolbar>
		</AppBar>
	);
}

export default TopBar;
