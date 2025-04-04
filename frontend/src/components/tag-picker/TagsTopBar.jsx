import { useTheme } from '@emotion/react';
import AddIcon from '@mui/icons-material/Add';
import MixIcon from '@mui/icons-material/AutoAwesome';
import ImageIcon from '@mui/icons-material/PermMedia';
import { Box, Divider, IconButton, Stack, Tooltip } from '@mui/material';
import React, { useState } from 'react';
import Client from '../../utils/client';
import TagsUtil from '../../utils/tags-util';
import ChooseDirectoryDialog from '../dialogs/ChooseDirectoryDialog';
import TextFieldWithKeyboard from '../text-field-with-keyboard/TextFieldWithKeyboard';
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
	onTitClicked,
	sortBy,
	setSortBy,
	prefixFilter,
	setPrefixFilter,
	tagSize,
	onZoomChanged,
	mixOnDemand,
}) {
	const theme = useTheme();
	const [showImagesFromDirectory, setShowImagesFromDirectory] = useState(false);

	const imageDirectoryChoosen = (directoryPath, doneCallback) => {
		Client.imageDirectoryChoosen(parentId, directoryPath).then(doneCallback);
	};

	return (
		<Stack flexDirection="column">
			<Box
				sx={{
					display: 'flex',
					flexDirection: 'row',
					padding: theme.spacing(1),
					gap: theme.spacing(1),
				}}
			>
				<Stack
					sx={{
						display: 'flex',
						width: '100%',
						flexDirection: 'row',
						gap: theme.spacing(1),
					}}
				>
					<Tooltip title="Generate one time mix">
						<IconButton onClick={() => mixOnDemand()}>
							<MixIcon sx={{ fontSize: theme.iconSize(1) }} />
						</IconButton>
					</Tooltip>
					{TagsUtil.allowToSetImageToCategory(parentId) && (
						<Tooltip title="Set images from directory">
							<IconButton onClick={() => setShowImagesFromDirectory(true)}>
								<ImageIcon sx={{ fontSize: theme.iconSize(1) }} />
							</IconButton>
						</Tooltip>
					)}
					{TagsUtil.allowToAddToCategory(parentId) && (
						<Tooltip title="Add new tag">
							<IconButton size="small" onClick={() => setAddTagDialogOpened(true)}>
								<AddIcon sx={{ fontSize: theme.iconSize(1) }} />
							</IconButton>
						</Tooltip>
					)}
					<TagImageTypeSelector
						disabled={false}
						tits={tits}
						tit={tit}
						onTitChanged={(newTit) => setTit(newTit)}
						onTitClicked={onTitClicked}
					/>
					<TextFieldWithKeyboard
						variant="outlined"
						autoFocus
						label="Search for tags..."
						type="search"
						size="small"
						onChange={(value) => setSearchTerm(value)}
					></TextFieldWithKeyboard>
					<TagSortSelector sortBy={sortBy} onSortChanged={(newSort) => setSortBy(newSort)} />
					<TagViewControls tagSize={tagSize} onZoomChanged={onZoomChanged} />
					<Divider orientation="vertical" />
					<PrefixFilter selectedChar={prefixFilter} setSelectedChar={setPrefixFilter} />
				</Stack>
			</Box>
			<TagsAnnotations
				annotations={annotations}
				selectedAnnotations={selectedAnnotations}
				setSelectedAnnotations={setSelectedAnnotations}
			/>
			<ChooseDirectoryDialog
				open={showImagesFromDirectory}
				title="Set images from directory"
				onChange={imageDirectoryChoosen}
				onClose={() => {
					setShowImagesFromDirectory(false);
				}}
			/>
		</Stack>
	);
}

export default TagsTopBar;
