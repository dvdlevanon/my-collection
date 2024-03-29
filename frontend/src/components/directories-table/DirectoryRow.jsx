import { useTheme } from '@emotion/react';
import { TableCell, TableRow } from '@mui/material';
import { useState } from 'react';
import { useQueryClient } from 'react-query';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import TimeUtil from '../../utils/time-utils';
import AddTagDialog from '../dialogs/AddTagDialog';
import DirectoryActionsCell from './DirectoryActionsCell';
import DirectoryCategoriesCell from './DirectoryCategoriesCell';
import DirectoryStatusCell from './DirectoryStatusCell';

function DirectoryRow({ directory }) {
	const queryClient = useQueryClient();
	const theme = useTheme();
	const [addCategoryDialogOpened, setAddCategoryDialogOpened] = useState(false);
	const refetchDirectories = () => {
		queryClient.refetchQueries({
			queryKey: ReactQueryUtil.DIRECTORIES_KEY,
		});
	};

	const excludeDirectory = (e, directory) => {
		Client.addOrUpdateDirectory({ ...directory, excluded: true }).then(refetchDirectories);
		ReactQueryUtil.updateDirectories(queryClient, directory.path, (currentDirectory) => {
			return {
				...currentDirectory,
				excluded: true,
			};
		});
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

	const setCategories = (categoryIds) => {
		let categories = [];
		for (let i = 0; i < categoryIds.length; i++) {
			categories.push({ id: categoryIds[i] });
		}

		Client.setDirectoryCategories({ ...directory, tags: categories }).then(refetchDirectories);
	};

	const formatLastSynced = (directory) => {
		if (!directory.lastSynced) {
			return 'Syncing...';
		}

		return TimeUtil.msToTime(Date.now() - directory.lastSynced) + ' Ago';
	};

	const formatFilesCount = (directory) => {
		if (directory.filesCount == undefined) {
			return 'N/A';
		}

		return directory.filesCount + ' files';
	};

	return (
		<>
			<TableRow
				sx={{
					opacity: directory.excluded ? '0.5' : '1',
				}}
				key={directory.path}
			>
				<TableCell>
					{!directory.excluded && <DirectoryStatusCell directory={directory} syncNow={syncNow} />}
				</TableCell>
				<TableCell>{directory.path}</TableCell>
				<TableCell>
					<DirectoryCategoriesCell
						directory={directory}
						setCategories={setCategories}
						onCreateCategoryClicked={() => setAddCategoryDialogOpened(true)}
					/>
				</TableCell>
				<TableCell>{!directory.excluded && formatFilesCount(directory)}</TableCell>
				<TableCell>{!directory.excluded && formatLastSynced(directory)}</TableCell>
				<TableCell>
					<DirectoryActionsCell
						directory={directory}
						includeDirectory={includeDirectory}
						excludeDirectory={excludeDirectory}
					/>
				</TableCell>
			</TableRow>
			<AddTagDialog
				open={addCategoryDialogOpened}
				parentId={null}
				verb="Category"
				onClose={() => setAddCategoryDialogOpened(false)}
			/>
		</>
	);
}

export default DirectoryRow;
