import CloseIcon from '@mui/icons-material/Close';
import { Dialog, DialogContent, DialogTitle, IconButton } from '@mui/material';
import React from 'react';
import { useQuery } from 'react-query';
import Client from '../network/client';
import ReactQueryUtil from '../utils/react-query-util';
import TagPicker from './TagPicker';

function AttachTagDialog({ open, item, onTagAdded, onClose }) {
	const tagsQuery = useQuery(ReactQueryUtil.TAGS_KEY, Client.getTags);

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
				<TagPicker
					initialSelectedSuperTagId={1}
					size="small"
					onTagSelected={onTagAdded}
					onDropDownToggled={() => {}}
				/>
			</DialogContent>
		</Dialog>
	);
}

export default AttachTagDialog;
