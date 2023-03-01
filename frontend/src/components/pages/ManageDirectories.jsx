import AddIcon from '@mui/icons-material/Add';
import { Fab, Typography } from '@mui/material';
import { Box } from '@mui/system';
import { useState } from 'react';
import { useQueryClient } from 'react-query';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import ChooseDirectoryDialog from '../dialogs/ChooseDirectoryDialog';
import DirectoriesTable from '../directories-table/DirectoriesTable';

function ManageDirectories() {
	const queryClient = useQueryClient();
	let [showAddDirectoryDialog, setShowAddDirectoryDialog] = useState(false);

	const addDirectory = (directoryPath, doneCallback) => {
		Client.addOrUpdateDirectory({ path: directoryPath }).then(() => {
			doneCallback();
		});
	};

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
				<ChooseDirectoryDialog
					title="Add Directory"
					onChange={addDirectory}
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
