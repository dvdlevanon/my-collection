import { useTheme } from '@emotion/react';
import CloseIcon from '@mui/icons-material/Close';
import { Button, Dialog, DialogActions, DialogContent, DialogTitle, IconButton, Typography } from '@mui/material';
import React from 'react';
import { useQueryClient } from '@tanstack/react-query';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';

function RemoveTagDialog({ tag, onClose }) {
	const queryClient = useQueryClient();
	const theme = useTheme();
	const removeTagClicked = (e) => {
		Client.removeTag(tag.id).then(() => {
			queryClient.refetchQueries({ queryKey: ReactQueryUtil.TAGS_KEY });
		});
	};

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
						top: '0',
						right: '0',
						margin: theme.spacing(1),
					}}
					onClick={(e) => {
						e.stopPropagation();
						onClose();
					}}
				>
					<CloseIcon sx={{ fontSize: theme.iconSize(1) }} />
				</IconButton>
			</DialogTitle>
			<DialogContent>
				<Typography variant="body1">
					Are you sure you want to remove {tag.title} with {tag.items ? tag.items.length : 0} items
				</Typography>
			</DialogContent>
			<DialogActions>
				<Button
					onClick={(e) => {
						e.stopPropagation();
						onClose();
					}}
				>
					Cancel
				</Button>
				<Button
					onClick={(e) => {
						e.stopPropagation();
						removeTagClicked();
					}}
				>
					Remove
				</Button>
			</DialogActions>
		</Dialog>
	);
}

export default RemoveTagDialog;
