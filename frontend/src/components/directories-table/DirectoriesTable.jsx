import { Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from '@mui/material';
import { useQuery, useQueryClient } from 'react-query';
import Client from '../../network/client';
import DirectoriesUtil from '../../utils/directories-util';
import ReactQueryUtil from '../../utils/react-query-util';
import DirectoryRow from './DirectoryRow';

function DirectoriesTable() {
	const queryClient = useQueryClient();
	const directoriesQuery = useQuery({
		queryKey: ReactQueryUtil.DIRECTORIES_KEY,
		queryFn: Client.getDirectories,
		onSuccess: (directories) => onDirectoriesSuccess(directories),
	});

	const onDirectoriesSuccess = (directories) => {
		if (directories.some(shouldRefetchDirectory)) {
			setTimeout(() => {
				queryClient.refetchQueries({
					queryKey: ReactQueryUtil.DIRECTORIES_KEY,
				});
			}, 1000);
		}
	};

	const shouldRefetchDirectory = (dir) => {
		return DirectoriesUtil.isProcessing(dir) && !DirectoriesUtil.isStaleProcessing(dir);
	};

	return (
		<TableContainer component={Paper}>
			<Table>
				<TableHead>
					<TableRow>
						<TableCell style={{ width: '5%' }}></TableCell>
						<TableCell style={{ width: '40%' }}>Path</TableCell>
						<TableCell style={{ width: '15%' }}>Category</TableCell>
						<TableCell style={{ width: '15%' }}>Files</TableCell>
						<TableCell style={{ width: '15%' }}>Last Scanned</TableCell>
						<TableCell style={{ width: '5%' }}></TableCell>
					</TableRow>
				</TableHead>
				<TableBody>
					{directoriesQuery.isSuccess &&
						directoriesQuery.data.map((directory) => {
							return <DirectoryRow key={directory.path} directory={directory} />;
						})}
				</TableBody>
			</Table>
		</TableContainer>
	);
}

export default DirectoriesTable;
