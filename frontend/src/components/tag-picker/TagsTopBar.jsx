import AddIcon from '@mui/icons-material/Add';
import ImageIcon from '@mui/icons-material/PermMedia';
import { Box, IconButton, TextField, Tooltip } from '@mui/material';
import React, { useState } from 'react';
import Client from '../../utils/client';
import TagsUtil from '../../utils/tags-util';
import ChooseDirectoryDialog from '../dialogs/ChooseDirectoryDialog';
import TagImageTypeSelector from './TagImageTypeSelector';
import TagSortSelector from './TagSortSelector';
import TagsAnnotations from './TagsAnnotations';

function TagsTopBar({
	parentId,
	annotations,
	selectedAnnotations,
	setSelectedAnnotations,
	setSearchTerm,
	setAddTagDialogOpened,
	tit,
	setTit,
	sortBy,
	setSortBy,
}) {
	let [showImagesFromDirectory, setShowImagesFromDirectory] = useState(false);
	const onSearchTermChanged = (e) => {
		setSearchTerm(e.target.value);
	};

	const imageDirectoryChoosen = (directoryPath, doneCallback) => {
		Client.imageDirectoryChoosen(parentId, directoryPath).then(doneCallback);
	};

	return (
		<Box
			sx={{
				display: 'flex',
				flexDirection: 'row',
				padding: '10px',
				gap: '10px',
			}}
		>
			{!TagsUtil.isDirectoriesCategory(parentId) && (
				<Tooltip title="Set images from directory">
					<IconButton onClick={() => setShowImagesFromDirectory(true)}>
						<ImageIcon />
					</IconButton>
				</Tooltip>
			)}
			{!TagsUtil.isDirectoriesCategory(parentId) && (
				<Tooltip title="Add new tag">
					<IconButton size="small" onClick={() => setAddTagDialogOpened(true)}>
						<AddIcon />
					</IconButton>
				</Tooltip>
			)}
			<TagImageTypeSelector tit={tit} onTitChanged={(newTit) => setTit(newTit)} />
			<TextField
				variant="outlined"
				autoFocus
				label="Search..."
				type="search"
				size="small"
				onChange={(e) => onSearchTermChanged(e)}
			></TextField>
			<TagSortSelector sortBy={sortBy} onSortChanged={(newSort) => setSortBy(newSort)} />
			<TagsAnnotations
				annotations={annotations}
				selectedAnnotations={selectedAnnotations}
				setSelectedAnnotations={setSelectedAnnotations}
			/>
			{showImagesFromDirectory && (
				<ChooseDirectoryDialog
					title="Set images from directory"
					onChange={imageDirectoryChoosen}
					onClose={() => {
						setShowImagesFromDirectory(false);
					}}
				/>
			)}
		</Box>
	);
}

export default TagsTopBar;
