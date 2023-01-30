import AddIcon from '@mui/icons-material/Add';
import ClearIcon from '@mui/icons-material/Clear';
import CloseIcon from '@mui/icons-material/Close';
import { Backdrop, Divider, IconButton, Popover, TextField, Typography } from '@mui/material';
import { Box } from '@mui/system';
import React, { useEffect, useRef, useState } from 'react';
import Client from '../network/client';

function TagAttachAnnotationMenu({ tag, menu, onClose }) {
	let [availableAnnotations, setAvailableAnnotations] = useState([]);
	const newAnnotationName = useRef(null);

	useEffect(
		() => Client.getAvailableAnnotations(tag.parentId, (annotations) => setAvailableAnnotations(annotations)),
		[]
	);

	const handleClose = (e) => {
		e.stopPropagation();
		onClose();
	};

	const addNewAnnotation = (e) => {
		e.stopPropagation();
		if (newAnnotationName.current.value == '') {
			return;
		}

		if (!tag.tags_annotations) {
			tag.tags_annotations = [];
		}

		tag.tags_annotations.push({ title: newAnnotationName.current.value });
		Client.saveTag(tag, () => {});
	};

	const removeAnnotation = (e, annotation) => {
		Client.removeAnnotationFromTag(tag.id, annotation.id, () => {});
		e.stopPropagation();
	};

	const addAnnotation = (e, annotation) => {
		Client.addAnnotationToTag(tag.id, annotation, () => {});
		e.stopPropagation();
	};

	const getAnnotations = (belongToTag) => {
		return availableAnnotations.filter((annotation) => {
			if (!tag.tags_annotations) {
				return !belongToTag;
			}

			if (tag.tags_annotations.some((cur) => cur.id == annotation.id)) {
				return belongToTag;
			} else {
				return !belongToTag;
			}
		});
	};

	return (
		<Popover
			slots={<Backdrop />}
			open={true}
			anchorReference="anchorPosition"
			anchorPosition={{ top: menu.mouseY, left: menu.mouseX }}
			BackdropProps={{
				invisible: false,
				onClick: (e) => {
					handleClose(e);
				},
			}}
			PaperProps={{
				sx: {
					display: 'flex',
					maxWidth: '400px',
					gap: '10px',
					flexDirection: 'column',
					padding: '10px',
				},
				onClick: (e) => e.stopPropagation(),
			}}
		>
			<Box>
				<Box
					onClick={(e) => e.stopPropagation()}
					sx={{
						display: 'flex',
						gap: '10px',
						alignItems: 'center',
					}}
				>
					<IconButton onClick={handleClose}>
						<CloseIcon />
					</IconButton>
					<Typography variant="h6" noWrap onClick={(e) => e.stopPropagation()}>
						{tag.title} Annotations
					</Typography>
				</Box>
				<Divider />
				<Box
					sx={{
						display: 'flex',
						flexDirection: 'row',
						flexWrap: 'wrap',
					}}
				>
					{getAnnotations(true).map((annotation) => {
						return (
							<Box
								bgcolor="primary.dark"
								color="primary.contrastText"
								sx={{
									margin: '10px',
									padding: '0px 10px',
									display: 'flex',
									cursor: 'pointer',
								}}
								borderRadius="10px"
								onClick={(e) => e.stopPropagation()}
								key={annotation.id}
							>
								<Typography sx={{ flexGrow: 1 }} variant="body1">
									{annotation.title}
								</Typography>
								<IconButton onClick={(e) => removeAnnotation(e, annotation)} size="small">
									<ClearIcon sx={{ fontSize: '15px' }} />
								</IconButton>
							</Box>
						);
					})}
					{getAnnotations(false).map((annotation) => {
						return (
							<Box
								bgcolor="gray"
								color="primary.contrastText"
								sx={{
									margin: '10px',
									padding: '0px 10px',
									display: 'flex',
									cursor: 'pointer',
								}}
								borderRadius="10px"
								onClick={(e) => addAnnotation(e, annotation)}
								key={annotation.id}
							>
								<Typography sx={{ flexGrow: 1 }} variant="body1">
									{annotation.title}
								</Typography>
							</Box>
						);
					})}
				</Box>
				<Box
					sx={{
						display: 'flex',
						gap: '10px',
						justifyContent: 'center',
						alignItems: 'center',
					}}
					onClick={(e) => {
						e.stopPropagation();
					}}
				>
					<TextField
						onKeyDown={(e) => {
							if (e.key == 'Enter') {
								addNewAnnotation(e);
							}
						}}
						autoFocus
						inputRef={newAnnotationName}
						placeholder="New Annotation Name..."
						sx={{
							flexGrow: 1,
						}}
					></TextField>
					<IconButton onClick={(e) => addNewAnnotation(e)} sx={{ alignSelf: 'center' }}>
						<AddIcon />
					</IconButton>
				</Box>
			</Box>
		</Popover>
	);
}

export default TagAttachAnnotationMenu;
