import AddIcon from '@mui/icons-material/Add';
import ImageIcon from '@mui/icons-material/PermMedia';
import { Box, IconButton, TextField, Tooltip } from '@mui/material';
import React, { useState } from 'react';
import Client from '../../utils/client';
import TagsUtil from '../../utils/tags-util';
import ChooseDirectoryDialog from '../dialogs/ChooseDirectoryDialog';
import TagAnnotation from './TagAnnotation';

function TagsFilter({
	parentId,
	annotations,
	selectedAnnotations,
	setSelectedAnnotations,
	setSearchTerm,
	setAddTagDialogOpened,
}) {
	let [showImagesFromDirectory, setShowImagesFromDirectory] = useState(false);
	const onSearchTermChanged = (e) => {
		setSearchTerm(e.target.value);
	};

	const isSelectedAnnotation = (annotation) => {
		return selectedAnnotations.some((cur) => annotation.id == cur.id);
	};

	const annotationSelected = (e, annotation) => {
		if (isSelectedAnnotation(annotation)) {
			setSelectedAnnotations(selectedAnnotations.filter((cur) => annotation.id != cur.id));
		} else {
			setSelectedAnnotations([...selectedAnnotations, annotation]);
		}
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
			<TextField
				variant="outlined"
				autoFocus
				label="Search..."
				type="search"
				size="small"
				onChange={(e) => onSearchTermChanged(e)}
			></TextField>
			<Box
				sx={{
					display: 'flex',
					flexDirection: 'row',
				}}
			>
				{annotations
					.sort((a, b) => (a.title > b.title ? 1 : a.title < b.title ? -1 : 0))
					.map((annotation) => {
						return (
							<TagAnnotation
								key={annotation.id}
								selectedAnnotaions
								annotation={annotation}
								selected={isSelectedAnnotation(annotation)}
								onClick={annotationSelected}
							/>
						);
					})}
			</Box>
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

export default TagsFilter;
