import CloseIcon from '@mui/icons-material/Close';
import { Box, Button, Dialog, DialogActions, DialogContent, DialogTitle, IconButton, TextField } from '@mui/material';
import React, { useRef } from 'react';

function ChooseDirectoryDialog({ title, onChange, onClose }) {
	const directoryFullPath = useRef(null);

	const addDirectory = () => {
		if (!directoryFullPath.current.value) {
			return;
		}

		onChange(directoryFullPath.current.value, () => {
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
			fullWidth={true}
			maxWidth={'sm'}
		>
			<DialogTitle variant="h6">
				{title}
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
			<DialogContent
				sx={{
					display: 'flex',
					gap: '10px',
				}}
			>
				<Box
					sx={{
						display: 'flex',
						flexDirection: 'column',
						gap: '10px',
						width: '100%',
					}}
				>
					<Box
						sx={{
							display: 'flex',
							flexDirection: 'row',
							gap: '10px',
						}}
					>
						<TextField
							sx={{
								flexGrow: 1,
							}}
							size="small"
							placeholder="Type directory full path here"
							inputRef={directoryFullPath}
							autoFocus
							focused
							onKeyDown={(e) => {
								if (e.key == 'Enter') {
									addDirectory(e);
								}
							}}
						></TextField>
					</Box>
				</Box>
			</DialogContent>
			<DialogActions>
				<Button color="secondary" onClick={onClose}>
					Cancel
				</Button>
				<Button variant="contained" onClick={(e) => addDirectory(e)}>
					Add
				</Button>
			</DialogActions>
		</Dialog>
	);
}

export default ChooseDirectoryDialog;
