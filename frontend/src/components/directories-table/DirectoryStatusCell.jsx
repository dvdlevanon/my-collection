import ErrorIcon from '@mui/icons-material/Error';
import SyncIcon from '@mui/icons-material/Sync';
import { CircularProgress, IconButton, Tooltip } from '@mui/material';
import DirectoriesUtil from '../../utils/directories-util';

function DirectoryStatusCell({ directory, syncNow }) {
	return (
		<>
			{!DirectoriesUtil.isProcessing(directory) && (
				<Tooltip title="Sync Now">
					<IconButton onClick={(e) => syncNow(e, directory)}>
						<SyncIcon color="secondary" />
					</IconButton>
				</Tooltip>
			)}
			{DirectoriesUtil.isProcessing(directory) && !DirectoriesUtil.isStaleProcessing(directory) && (
				<Tooltip title="Processing">
					<IconButton>
						<CircularProgress color="secondary" size="25px" />
					</IconButton>
				</Tooltip>
			)}
			{DirectoriesUtil.isProcessing(directory) && DirectoriesUtil.isStaleProcessing(directory) && (
				<Tooltip title="Error processing, click to try again">
					<IconButton onClick={(e) => syncNow(e, directory)}>
						<ErrorIcon color="error" />
					</IconButton>
				</Tooltip>
			)}
		</>
	);
}

export default DirectoryStatusCell;
