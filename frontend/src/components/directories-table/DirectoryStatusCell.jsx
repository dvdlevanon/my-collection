import { useTheme } from '@emotion/react';
import ErrorIcon from '@mui/icons-material/Error';
import SyncIcon from '@mui/icons-material/Sync';
import { CircularProgress, IconButton, Tooltip } from '@mui/material';
import DirectoriesUtil from '../../utils/directories-util';

function DirectoryStatusCell({ directory, syncNow }) {
	const theme = useTheme();

	return (
		<>
			{!DirectoriesUtil.isProcessing(directory) && (
				<Tooltip title="Sync Now">
					<IconButton onClick={(e) => syncNow(e, directory)}>
						<SyncIcon color="secondary" sx={{ fontSize: theme.iconSize(1) }} />
					</IconButton>
				</Tooltip>
			)}
			{DirectoriesUtil.isProcessing(directory) && !DirectoriesUtil.isStaleProcessing(directory) && (
				<Tooltip title="Processing">
					<IconButton>
						<CircularProgress color="secondary" size={theme.iconSize(1)} />
					</IconButton>
				</Tooltip>
			)}
			{DirectoriesUtil.isProcessing(directory) && DirectoriesUtil.isStaleProcessing(directory) && (
				<Tooltip title="Error processing, click to try again">
					<IconButton onClick={(e) => syncNow(e, directory)}>
						<ErrorIcon color="error" sx={{ fontSize: theme.iconSize(1) }} />
					</IconButton>
				</Tooltip>
			)}
		</>
	);
}

export default DirectoryStatusCell;
