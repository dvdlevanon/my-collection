import CloseIcon from '@mui/icons-material/Close';
import { Dialog, DialogContent, DialogTitle, IconButton } from '@mui/material';
import React from 'react';
import TagChooser from './TagChooser';

function AttachTagDialog({ open, item, tags, onTagAdded, onClose }) {
	return (
		<Dialog
			onClose={(e, reason) => {
				if (reason == 'backdropClick' || reason == 'escapeKeyDown') {
					onClose(e);
				}
			}}
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
			<DialogContent>
				<TagChooser
					initialSelectedSuperTag={tags[0]}
					tags={tags}
					size="small"
					onTagSelected={onTagAdded}
					onDropDownToggled={() => {}}
				/>
			</DialogContent>
		</Dialog>
	);
}

export default AttachTagDialog;
