import AddIcon from '@mui/icons-material/Add';
import { Fab, Typography } from '@mui/material';
import { Box } from '@mui/system';
import { useState } from 'react';
import { useQueryClient } from 'react-query';
import ReactQueryUtil from '../../utils/react-query-util';
import AddDirectoryDialog from '../dialogs/AddDirectoryDialog';
import DirectoriesTable from '../directories-table/DirectoriesTable';

function ManageDirectories() {
	const queryClient = useQueryClient();
	let [rootDirectory, setRootDirectory] = useState('/mnt/usb1');
	let [showAddDirectoryDialog, setShowAddDirectoryDialog] = useState(false);

	return (
		<Box
			sx={{
				height: '95%',
				padding: '20px',
				position: 'relative',
				display: 'flex',
				flexDirection: 'column',
			}}
		>
			<Typography
				sx={{
					padding: '0px 0px 15px 0px',
				}}
				variant="h6"
			>
				Manage Directories
			</Typography>
			<Box
				sx={{
					flexGrow: 1,
					overflowX: 'hidden',
					overflowY: 'scroll',
				}}
			>
				<DirectoriesTable />
			</Box>
			<Fab
				sx={{
					position: 'absolute',
					bottom: '30px',
					right: '50px',
				}}
				color="primary"
				onClick={() => setShowAddDirectoryDialog(true)}
			>
				<AddIcon />
			</Fab>
			{showAddDirectoryDialog && (
				<AddDirectoryDialog
					onClose={() => {
						setShowAddDirectoryDialog(false);
						queryClient.refetchQueries({ queryKey: ReactQueryUtil.DIRECTORIES_KEY });
					}}
				/>
			)}
		</Box>
	);
}

export default ManageDirectories;
