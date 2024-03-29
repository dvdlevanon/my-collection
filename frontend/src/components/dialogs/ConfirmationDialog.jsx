import { useTheme } from '@emotion/react';
import CloseIcon from '@mui/icons-material/Close';
import { Button, Dialog, DialogActions, DialogContent, DialogTitle, IconButton, Typography } from '@mui/material';
import React from 'react';

function ConfirmationDialog({ open, title, text, actionButtonTitle, onConfirm, onCancel }) {
	const theme = useTheme();

	return (
		<Dialog
			open={open}
			onClose={(e, reason) => {
				e.stopPropagation();
				if (reason == 'backdropClick' || reason == 'escapeKeyDown') {
					onCancel();
				}
			}}
		>
			<DialogTitle>
				{title}
				<IconButton
					sx={{
						position: 'absolute',
						top: '0',
						right: '0',
						margin: theme.spacing(1),
					}}
					onClick={(e) => {
						e.stopPropagation();
						onCancel();
					}}
				>
					<CloseIcon sx={{ fontSize: theme.iconSize(1) }} />
				</IconButton>
			</DialogTitle>
			<DialogContent>
				<Typography variant="body1">{text}</Typography>
			</DialogContent>
			<DialogActions>
				<Button
					onClick={(e) => {
						e.stopPropagation();
						onCancel();
					}}
				>
					Cancel
				</Button>
				<Button
					onClick={(e) => {
						e.stopPropagation();
						onConfirm();
					}}
				>
					{actionButtonTitle}
				</Button>
			</DialogActions>
		</Dialog>
	);
}

export default ConfirmationDialog;
