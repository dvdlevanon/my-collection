import { Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from '@mui/material';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { useEffect } from 'react';
import Client from '../../utils/client';
import DirectoriesUtil from '../../utils/directories-util';
import ReactQueryUtil from '../../utils/react-query-util';
import DirectoryRow from './DirectoryRow';

function DirectoriesTable() {
	const queryClient = useQueryClient();
	const directoriesQuery = useQuery({
		queryKey: ReactQueryUtil.DIRECTORIES_KEY,
		queryFn: Client.getDirectories,
	});

	useEffect(() => {
		if (directoriesQuery.data) {
			onDirectoriesSuccess(directoriesQuery.data);
		}
	}, [directoriesQuery.data]);

	const onDirectoriesSuccess = (directories) => {
		if (directories.some(shouldRefetchDirectory)) {
			setTimeout(() => {
				queryClient.refetchQueries({
					queryKey: ReactQueryUtil.DIRECTORIES_KEY,
				});
			}, 5000);
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
