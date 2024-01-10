import { useTheme } from '@emotion/react';
import CloseIcon from '@mui/icons-material/Close';
import { Dialog, DialogContent, DialogTitle, IconButton } from '@mui/material';
import React from 'react';
import { useQuery } from 'react-query';
import Client from '../../utils/client';
import GalleryUrlParams from '../../utils/gallery-url-params';
import ReactQueryUtil from '../../utils/react-query-util';
import TagPicker from '../tag-picker/TagPicker';

function AttachTagDialog({ open, item, onTagAdded, onClose, singleCategoryMode, initialSelectedCategoryId }) {
	const tagsQuery = useQuery(ReactQueryUtil.TAGS_KEY, Client.getTags);
	const theme = useTheme();

	return (
		<Dialog
			onClose={(e, reason) => {
				if (reason == 'backdropClick' || reason == 'escapeKeyDown') {
					onClose(e);
				}
			}}
			disableEnforceFocus
			open={open}
			fullWidth={true}
			maxWidth={'xl'}
			PaperProps={{
				sx: {
					height: '70%',
				},
			}}
		>
			<DialogTitle variant="h6">
				Add a tag to {item.title}
				<IconButton
					sx={{
						position: 'absolute',
						top: '0',
						right: '0',
						margin: theme.spacing(1),
					}}
					onClick={onClose}
				>
					<CloseIcon sx={{ fontSize: theme.iconSize(1) }} />
				</IconButton>
			</DialogTitle>
			<DialogContent sx={{ overflow: 'hidden' }}>
				<TagPicker
					origin="attach-dialog"
					initialSelectedCategoryId={initialSelectedCategoryId}
					singleCategoryMode={singleCategoryMode}
					showSpecialCategories={false}
					initialTagSize={350}
					onTagSelected={onTagAdded}
					onDropDownToggled={() => {}}
					tagLinkBuilder={(tag) => '/?' + GalleryUrlParams.buildUrlParams(tag.id)}
					setHideTopBar={() => {}}
				/>
			</DialogContent>
		</Dialog>
	);
}

export default AttachTagDialog;
