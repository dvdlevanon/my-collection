import AddIcon from '@mui/icons-material/Add';
import CloseIcon from '@mui/icons-material/Close';
import { IconButton, Menu, MenuItem, TextField, Typography } from '@mui/material';
import { Box } from '@mui/system';
import React, { useEffect, useRef, useState } from 'react';
import Client from '../network/client';

function TagAttachAnnotationMenu({ tag, menu, onClose }) {
	let [availableAnnotations, setAvailableAnnotations] = useState([]);
	const newAttributeName = useRef(null);

	useEffect(
		() => Client.getAvailableAnnotations(tag.parentId, (annotations) => setAvailableAnnotations(annotations)),
		[]
	);

	const handleClose = (e) => {
		e.stopPropagation();
		onClose();
	};

	const addNewAttributeClicked = (e) => {
		if (!tag.tags_annotations) {
			tag.tags_annotations = [];
		}

		tag.tags_annotations.push({ title: newAttributeName.current.value });
		Client.saveTag(tag, () => {});
	};

	return (
		<Menu open={true} anchorReference="anchorPosition" anchorPosition={{ top: menu.mouseY, left: menu.mouseX }}>
			<MenuItem onClick={(e) => e.stopPropagation()}>
				<Box sx={{ display: 'flex', gap: '10px', alignItems: 'center' }}>
					<IconButton onClick={handleClose}>
						<CloseIcon />
					</IconButton>
					<Typography onClick={(e) => e.stopPropagation()}>Attach Attribute to {tag.title}</Typography>
				</Box>
			</MenuItem>
			<MenuItem onClick={(e) => handleClose(e)}>
				<Box
					sx={{ display: 'flex', gap: '10px', justifyContent: 'center', alignItems: 'center' }}
					onClick={(e) => {
						e.stopPropagation();
					}}
				>
					<TextField autoFocus inputRef={newAttributeName} placeholder="New Attribute Name..."></TextField>
					<IconButton sx={{ alignSelf: 'center' }}>
						<AddIcon onClick={(e) => addNewAttributeClicked(e)} />
					</IconButton>
				</Box>
			</MenuItem>
		</Menu>
	);
}

export default TagAttachAnnotationMenu;
