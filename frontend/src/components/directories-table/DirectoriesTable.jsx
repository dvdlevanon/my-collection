import AddIcon from '@mui/icons-material/Add';
import DeleteIcon from '@mui/icons-material/Delete';
import ErrorIcon from '@mui/icons-material/Error';
import SyncIcon from '@mui/icons-material/Sync';
import {
	CircularProgress,
	IconButton,
	Paper,
	Table,
	TableBody,
	TableCell,
	TableContainer,
	TableHead,
	TableRow,
	Tooltip,
} from '@mui/material';
import { useQuery, useQueryClient } from 'react-query';
import Client from '../../network/client';
import ReactQueryUtil from '../../utils/react-query-util';
function DirectoriesTable() {
	const queryClient = useQueryClient();
	const directoriesQuery = useQuery({
		queryKey: ReactQueryUtil.DIRECTORIES_KEY,
		queryFn: Client.getDirectories,
		onSuccess: (directories) => {
			if (directories.some((dir) => isProcessing(dir) && !isStaleProcessing(dir))) {
				console.log('Setting timeout');
				setTimeout(refetchDirectories, 1000);
			}
		},
	});

	const refetchDirectories = () => {
		console.log('Refetching');
		queryClient.refetchQueries({
			queryKey: ReactQueryUtil.DIRECTORIES_KEY,
		});
	};

	const removeDirectory = (e, directory) => {
		Client.removeDirectory(directory.path).then(refetchDirectories);
	};

	const includeDirectory = (e, directory) => {
		Client.addOrUpdateDirectory({ ...directory, excluded: false }).then(refetchDirectories);
		ReactQueryUtil.updateDirectories(queryClient, directory.path, (currentDirectory) => {
			return {
				...currentDirectory,
				excluded: false,
				processingStart: Date.now(),
			};
		});
	};

	const syncNow = (e, directory) => {
		Client.addOrUpdateDirectory(directory).then(refetchDirectories);
		ReactQueryUtil.updateDirectories(queryClient, directory.path, (currentDirectory) => {
			return {
				...currentDirectory,
				processingStart: Date.now(),
			};
		});
	};

	const msToTime = (millis) => {
		let seconds = (millis / 1000).toFixed(1);
		let minutes = (millis / (1000 * 60)).toFixed(1);
		let hours = (millis / (1000 * 60 * 60)).toFixed(1);
		let days = (millis / (1000 * 60 * 60 * 24)).toFixed(1);
		if (seconds < 60) return Math.floor(seconds) + ' Seconds';
		else if (minutes < 60) return Math.floor(minutes) + ' Minutes';
		else if (hours < 24) return Math.floor(hours) + ' Hours';
		else return Math.floor(days) + ' Days';
	};

	const formatLastSynced = (directory) => {
		if (!directory.lastSynced) {
			return 'Syncing...';
		}

		return msToTime(Date.now() - directory.lastSynced) + ' Ago';
	};

	const formatFilesCount = (directory) => {
		if (directory.filesCount == undefined) {
			return 'N/A';
		}

		return directory.filesCount + ' files';
	};

	const isProcessing = (directory) => {
		return directory.processingStart != undefined && directory.processingStart > 0;
	};

	const isStaleProcessing = (directory) => {
		if (!isProcessing(directory)) {
			return false;
		}

		let millisSinceStart = Date.now() - directory.processingStart;
		return millisSinceStart > 1000 * 60;
	};

	return (
		<TableContainer component={Paper}>
			<Table>
				<TableHead>
					<TableRow>
						<TableCell style={{ width: '5%' }}></TableCell>
						<TableCell style={{ width: '20%' }}>Path</TableCell>
						<TableCell style={{ width: '20%' }}>Files</TableCell>
						<TableCell style={{ width: '20%' }}>Automatic Tags</TableCell>
						<TableCell style={{ width: '20%' }}>Last Scanned</TableCell>
						<TableCell style={{ width: '5%' }}></TableCell>
					</TableRow>
				</TableHead>
				<TableBody>
					{directoriesQuery.isSuccess &&
						directoriesQuery.data.map((directory) => {
							return (
								<TableRow
									sx={{
										backgroundColor: directory.excluded ? '#333' : 'main',
									}}
									key={directory.path}
								>
									<TableCell>
										{!directory.excluded && !isProcessing(directory) && (
											<Tooltip title="Sync Now">
												<IconButton onClick={(e) => syncNow(e, directory)}>
													<SyncIcon color="secondary" />
												</IconButton>
											</Tooltip>
										)}
										{!directory.excluded &&
											isProcessing(directory) &&
											!isStaleProcessing(directory) && (
												<Tooltip title="Processing">
													<IconButton>
														<CircularProgress color="secondary" size="25px" />
													</IconButton>
												</Tooltip>
											)}
										{!directory.excluded &&
											isProcessing(directory) &&
											isStaleProcessing(directory) && (
												<Tooltip title="Error processing">
													<IconButton>
														<ErrorIcon color="error" />
													</IconButton>
												</Tooltip>
											)}
									</TableCell>
									<TableCell>{directory.path}</TableCell>
									<TableCell>{!directory.excluded && formatFilesCount(directory)}</TableCell>
									<TableCell>
										{!directory.excluded && directory.tags && directory.tags.join('\n')}
									</TableCell>
									<TableCell>{!directory.excluded && formatLastSynced(directory)}</TableCell>
									<TableCell>
										{!directory.excluded && (
											<Tooltip title="Delete">
												<IconButton onClick={(e) => removeDirectory(e, directory)}>
													<DeleteIcon color="secondary" />
												</IconButton>
											</Tooltip>
										)}
										{directory.excluded && (
											<Tooltip title="Add">
												<IconButton onClick={(e) => includeDirectory(e, directory)}>
													<AddIcon color="secondary" />
												</IconButton>
											</Tooltip>
										)}
									</TableCell>
								</TableRow>
							);
						})}
				</TableBody>
			</Table>
		</TableContainer>
	);
}

export default DirectoriesTable;
