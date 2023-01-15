import styles from './AddTagDialog.module.css';
import React from 'react';
import { Dialog, DialogContent, DialogTitle, IconButton } from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import TagChooser from './TagChooser';

function AddTagDialog({ open, item, tags, onTagAdded, onClose }) {
	return (
		<Dialog open={open} fullWidth maxWidth="false">
			<DialogTitle>
				Add a tag to {item.title}
				<IconButton className={styles.close_button} onClick={onClose}>
					<CloseIcon />
				</IconButton>
			</DialogTitle>
			<DialogContent className={styles.dialog_content}>
				<TagChooser tags={tags} onTagSelected={onTagAdded} />
			</DialogContent>
		</Dialog>
	);
}

export default AddTagDialog;
