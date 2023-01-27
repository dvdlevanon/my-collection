import CloseIcon from '@mui/icons-material/Close';
import { Dialog, DialogContent, DialogTitle, IconButton, Typography } from '@mui/material';
import React from 'react';
import styles from './AddTagDialog.module.css';
import TagChooser from './TagChooser';

function AddTagDialog({ open, item, tags, onTagAdded, onClose }) {
	return (
		<Dialog open={open} modal={true} fullWidth={true} maxWidth={'md'}>
			<DialogTitle>
				<Typography variant="caption">Add a tag to {item.title}</Typography>
				<IconButton className={styles.close_button} onClick={onClose}>
					<CloseIcon />
				</IconButton>
			</DialogTitle>
			<DialogContent className={styles.dialog_content}>
				<TagChooser tags={tags} size="small" onTagSelected={onTagAdded} onDropDownToggled={() => {}} />
			</DialogContent>
		</Dialog>
	);
}

export default AddTagDialog;
