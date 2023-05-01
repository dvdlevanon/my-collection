import AddIcon from '@mui/icons-material/Add';
import ImageIcon from '@mui/icons-material/PermMedia';
import { Box, IconButton, Stack, TextField, Tooltip } from '@mui/material';
import React, { useState } from 'react';
import Client from '../../utils/client';
import TagsUtil from '../../utils/tags-util';
import ChooseDirectoryDialog from '../dialogs/ChooseDirectoryDialog';
import PrefixFilter from './PrefixFilter';
import TagImageTypeSelector from './TagImageTypeSelector';
import TagSortSelector from './TagSortSelector';
import TagViewControls from './TagViewControls';
import TagsAnnotations from './TagsAnnotations';

function TagsTopBar({
	parentId,
	annotations,
	selectedAnnotations,
	setSelectedAnnotations,
	setSearchTerm,
	setAddTagDialogOpened,
	tits,
	tit,
	setTit,
	sortBy,
	setSortBy,
	prefixFilter,
	setPrefixFilter,
	tagSize,
	onZoomChanged,
}) {
	let [showImagesFromDirectory, setShowImagesFromDirectory] = useState(false);

	const imageDirectoryChoosen = (directoryPath, doneCallback) => {
		Client.imageDirectoryChoosen(parentId, directoryPath).then(doneCallback);
	};

	return (
		<Stack flexDirection="column">
			<Box
				sx={{
					display: 'flex',
					flexDirection: 'row',
					padding: '10px',
					gap: '10px',
				}}
			>
				<Stack
					sx={{
						display: 'flex',
						flexDirection: 'row',
						gap: '10px',
						maxHeight: '44px',
					}}
				>
					{TagsUtil.allowToSetImageToCategory(parentId) && (
						<Tooltip title="Set images from directory">
							<IconButton onClick={() => setShowImagesFromDirectory(true)}>
								<ImageIcon />
							</IconButton>
						</Tooltip>
					)}
					{TagsUtil.allowToAddToCategory(parentId) && (
						<Tooltip title="Add new tag">
							<IconButton size="small" onClick={() => setAddTagDialogOpened(true)}>
								<AddIcon />
							</IconButton>
						</Tooltip>
					)}
					<TagImageTypeSelector tits={tits} tit={tit} onTitChanged={(newTit) => setTit(newTit)} />
					<TextField
						variant="outlined"
						autoFocus
						label="Search for tags..."
						type="search"
						size="small"
						onChange={(e) => setSearchTerm(e.target.value)}
						sx={{
							minWidth: '200px',
						}}
					></TextField>
					<TagSortSelector sortBy={sortBy} onSortChanged={(newSort) => setSortBy(newSort)} />
					<TagViewControls tagSize={tagSize} onZoomChanged={onZoomChanged} />
				</Stack>
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
			<PrefixFilter selectedChar={prefixFilter} setSelectedChar={setPrefixFilter} />
		</Stack>
	);
}

export default TagsTopBar;
