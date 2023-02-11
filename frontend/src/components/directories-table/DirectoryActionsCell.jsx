import AddIcon from '@mui/icons-material/Add';
import DeleteIcon from '@mui/icons-material/Delete';
import { IconButton, Tooltip } from '@mui/material';
import React from 'react';

function DirectoryActionsCell({ directory, includeDirectory, excludeDirectory }) {
	return (
		<>
			{(!directory.excluded && (
				<Tooltip title="Delete">
					<IconButton onClick={(e) => excludeDirectory(e, directory)}>
						<DeleteIcon color="secondary" />
					</IconButton>
				</Tooltip>
			)) || (
				<Tooltip title="Add">
					<IconButton onClick={(e) => includeDirectory(e, directory)}>
						<AddIcon color="secondary" />
					</IconButton>
				</Tooltip>
			)}
		</>
	);
}

export default DirectoryActionsCell;
