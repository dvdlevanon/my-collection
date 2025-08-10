import { useTheme } from '@emotion/react';
import AddIcon from '@mui/icons-material/Add';
import { Fab, Stack, Typography } from '@mui/material';
import { Box } from '@mui/system';
import { useState } from 'react';
import { useQueryClient } from 'react-query';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import ChooseDirectoryDialog from '../dialogs/ChooseDirectoryDialog';
import DirectoriesTable from '../directories-table/DirectoriesTable';

function ManageDirectories() {
	const queryClient = useQueryClient();
	const [showAddDirectoryDialog, setShowAddDirectoryDialog] = useState(false);
	const theme = useTheme();

	const addDirectory = (directoryPath, doneCallback) => {
		Client.addOrUpdateDirectory({ path: directoryPath }).then(() => {
			doneCallback();
		});
	};

	return (
		<Stack height="95%" padding={theme.spacing(1)} position="relative" flexDirection="column">
			<Typography padding={theme.multiSpacing(0, 0, 1.5, 0)} variant="h6">
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
					bottom: theme.spacing(3),
					right: theme.spacing(5),
				}}
				color="primary"
				onClick={() => setShowAddDirectoryDialog(true)}
			>
				<AddIcon sx={{ fontSize: theme.iconSize(1) }} />
			</Fab>
			{showAddDirectoryDialog && (
				<ChooseDirectoryDialog
					open={true}
					title="Add Directory"
					onChange={addDirectory}
					onClose={() => {
						setShowAddDirectoryDialog(false);
						queryClient.refetchQueries({ queryKey: ReactQueryUtil.DIRECTORIES_KEY });
					}}
				/>
			)}
		</Stack>
	);
}

export default ManageDirectories;
