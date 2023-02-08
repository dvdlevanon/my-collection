import AddIcon from '@mui/icons-material/Add';
import DeleteIcon from '@mui/icons-material/Delete';
import SyncIcon from '@mui/icons-material/Sync';
import {
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
	const directoriesQuery = useQuery(ReactQueryUtil.DIRECTORIES_KEY, Client.getDirectories);

	const refetchDirectories = () => {
		queryClient.refetchQueries({
			queryKey: ReactQueryUtil.DIRECTORIES_KEY,
		});
	};

	const removeDirectory = (e, directory) => {
		Client.removeDirectory(directory.path).then(refetchDirectories);
	};

	const includeDirectory = (e, directory) => {
		Client.addOrUpdateDirectory({ ...directory, excluded: false }).then(refetchDirectories);
	};

	const syncNow = (e, directory) => {
		Client.addOrUpdateDirectory({ directory }).then(refetchDirectories);
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
		return msToTime(Date.now() - directory.lastSynced) + ' Ago';
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
										{!directory.excluded && (
											<Tooltip title="Sync Now">
												<IconButton onClick={(e) => syncNow(e, directory)}>
													<SyncIcon color="secondary" />
												</IconButton>
											</Tooltip>
										)}
									</TableCell>
									<TableCell>{directory.path}</TableCell>
									<TableCell>{!directory.excluded && directory.filesCount + ' files'}</TableCell>
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
