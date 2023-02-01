import CloseIcon from '@mui/icons-material/Close';
import { Button, Dialog, DialogActions, DialogContent, DialogTitle, IconButton, Typography } from '@mui/material';
import React from 'react';

function RemoveTag({ tag, onClose }) {
	return (
		<Dialog
			onClose={(e, reason) => {
				e.stopPropagation();
				if (reason == 'backdropClick' || reason == 'escapeKeyDown') {
					onClose();
				}
			}}
			open={true}
		>
			<DialogTitle>
				Remove Tag
				<IconButton
					sx={{
						position: 'absolute',
						top: '0px',
						right: '0px',
						margin: '10px',
					}}
					onClick={(e) => {
						e.stopPropagation();
						onClose();
					}}
				>
					<CloseIcon />
				</IconButton>
			</DialogTitle>
			<DialogContent>
				<Typography variant="body1">
					Are you sure you want to remove {tag.title} with {tag.items ? tag.items.length : 0} items
				</Typography>
			</DialogContent>
			<DialogActions>
				<Button>Cancel</Button>
			</DialogActions>
		</Dialog>
	);
}

export default RemoveTag;
