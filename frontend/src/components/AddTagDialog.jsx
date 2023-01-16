import CloseIcon from '@mui/icons-material/Close';
import { Dialog, DialogContent, DialogTitle, IconButton } from '@mui/material';
import React from 'react';
import styles from './AddTagDialog.module.css';
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
				<TagChooser tags={tags} markActive={false} onTagSelected={onTagAdded} />
			</DialogContent>
		</Dialog>
	);
}

export default AddTagDialog;
