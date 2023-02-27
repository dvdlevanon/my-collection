import { TableCell, TableRow } from '@mui/material';
import { useState } from 'react';
import { useQueryClient } from 'react-query';
import Client from '../../network/client';
import ReactQueryUtil from '../../utils/react-query-util';
import AddTagDialog from '../dialogs/AddTagDialog';
import DirectoryActionsCell from './DirectoryActionsCell';
import DirectoryCategoriesCell from './DirectoryCategoriesCell';
import DirectoryStatusCell from './DirectoryStatusCell';

function DirectoryRow({ directory }) {
	const queryClient = useQueryClient();
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

	return (
		<>
			<TableRow
				sx={{
					backgroundColor: directory.excluded ? '#333' : 'main',
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
			{addCategoryDialogOpened && (
				<AddTagDialog parentId={null} verb="Category" onClose={() => setAddCategoryDialogOpened(false)} />
			)}
		</>
	);
}

export default DirectoryRow;
