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
						top: '0px',
						right: '0px',
						margin: '10px',
					}}
					onClick={onClose}
				>
					<CloseIcon />
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
