import AddIcon from '@mui/icons-material/Add';
import CloseIcon from '@mui/icons-material/Close';
import { Dialog, DialogContent, DialogTitle, IconButton, TextField } from '@mui/material';
import React, { useRef } from 'react';
import { useQueryClient } from 'react-query';
import Client from '../../network/client';
import ReactQueryUtil from '../../utils/react-query-util';

function AddTagDialog({ parentId, onClose }) {
	const queryClient = useQueryClient();
	const newTagName = useRef(null);

	const addTag = (e) => {
		if (!newTagName.current.value) {
			return;
		}

		let newTag = {
			title: newTagName.current.value,
			parentId: parentId,
		};

		Client.createTag(newTag)
			.then((response) => response.json())
			.then((newTag) => {
				queryClient.refetchQueries({ queryKey: ReactQueryUtil.TAGS_KEY });
				onClose();
			});
	};

	return (
		<Dialog
			onClose={(e, reason) => {
				if (reason == 'backdropClick' || reason == 'escapeKeyDown') {
					onClose();
				}
			}}
			open={true}
		>
			<DialogTitle variant="h6">
				Add Tag
				<IconButton
					sx={{
						position: 'absolute',
						top: '0px',
						right: '0px',
						margin: '10px',
					}}
					onClick={() => onClose()}
				>
					<CloseIcon />
				</IconButton>
			</DialogTitle>
			<DialogContent
				sx={{
					display: 'flex',
					gap: '10px',
				}}
			>
				<TextField
					autoFocus
					onKeyDown={(e) => {
						if (e.key == 'Enter') {
							addTag(e);
						}
					}}
					size="small"
					placeholder="Tag Name"
					inputRef={newTagName}
				></TextField>
				<IconButton onClick={(e) => addTag(e)} sx={{ alignSelf: 'center' }}>
					<AddIcon />
				</IconButton>
			</DialogContent>
		</Dialog>
	);
}

export default AddTagDialog;
