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
	Stack,
	Toolbar,
	Tooltip,
	Typography,
} from '@mui/material';
import { useQuery } from '@tanstack/react-query';
import { useEffect, useState } from 'react';
import { Link as RouterLink } from 'react-router-dom';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import TimeUtil from '../../utils/time-utils';
import Queue from '../queue/Queue';
import ThemeSelector from '../theme-selector/ThemeSelector';

function TopBar({ previewMode, onPreviewModeChange, theme, setTheme }) {
	const [refreshAnchorEl, setRefreshAnchorEl] = useState(null);
	const queueMetadataQuery = useQuery({
		queryKey: ReactQueryUtil.QUEUE_METADATA_KEY,
		queryFn: Client.getQueueMetadata,
	});

	const statsQuery = useQuery({
		queryKey: ReactQueryUtil.STATS_KEY,
		queryFn: Client.getStats,
	});

	useEffect(() => {
		if (statsQuery.data) {
			if (statsQuery.data.size == 0) {
				setQueueEl(null);
			}
		}
	}, [statsQuery.data]);

	const previewsChange = (e) => {
		onPreviewModeChange(e.target.checked);
	};

	const [queueEl, setQueueEl] = useState(null);

	return (
		<AppBar position="static">
			<Toolbar
				sx={{
					gap: theme.spacing(2),
					alignItems: 'center',
					alignContent: 'center',
					verticalAlign: 'center',
				}}
			>
				<Link color="inherit" sx={{ flexGrow: 1 }} component={RouterLink} to={'/' + window.location.search}>
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
								<Badge badgeContent={queueMetadataQuery.data.unfinishedTasks || null}>
									<QueueIcon sx={{ fontSize: theme.iconSize(1) }} />
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
					<Button variant="contained">Manage Directories</Button>
				</Link>
				<Link href={Client.getExportMetadataUrl()} download>
					<Button variant="contained">Export metadata</Button>
				</Link>
				<Button variant="contained" onClick={(e) => setRefreshAnchorEl(e.currentTarget)}>
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
					<MenuItem
						onClick={() => {
							Client.refreshFileMetadata();
							setRefreshAnchorEl(null);
						}}
					>
						Refresh File Metadata
					</MenuItem>
					<MenuItem
						onClick={() => {
							Client.runSpectagger();
							setRefreshAnchorEl(null);
						}}
					>
						Run spectagger
					</MenuItem>
					<MenuItem
						onClick={() => {
							Client.runItemOptimizer();
							setRefreshAnchorEl(null);
						}}
					>
						Run item optimizer
					</MenuItem>
					<MenuItem
						onClick={() => {
							Client.runDirectoryScan();
							setRefreshAnchorEl(null);
						}}
					>
						Run directory scan
					</MenuItem>
				</Menu>
				<FormControlLabel
					label="Use Previews"
					control={<Checkbox checked={previewMode} onChange={(e) => previewsChange(e)} />}
				/>
				{statsQuery.isSuccess && (
					<Stack flexDirection="row" gap={theme.spacing(1)} justifyContent="center" alignItems="center">
						<Stack flexDirection="column">
							<Typography variant="caption">Tags: {statsQuery.data.tags_count}</Typography>
							<Typography variant="caption">Items: {statsQuery.data.items_count}</Typography>
						</Stack>
						<Typography variant="caption">
							Total Duration: {TimeUtil.msToTime(statsQuery.data.total_duration_seconds * 1000)}{' '}
						</Typography>
					</Stack>
				)}
				<ThemeSelector theme={theme} setTheme={setTheme} />
			</Toolbar>
		</AppBar>
	);
}

export default TopBar;
