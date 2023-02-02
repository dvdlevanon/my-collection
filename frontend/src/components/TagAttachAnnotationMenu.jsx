import AddIcon from '@mui/icons-material/Add';
import CloseIcon from '@mui/icons-material/Close';
import { Divider, IconButton, Popover, TextField, Typography } from '@mui/material';
import { Box } from '@mui/system';
import React, { useRef } from 'react';
import { useMutation, useQuery, useQueryClient } from 'react-query';
import Client from '../network/client';
import ReactQueryUtil from '../utils/react-query-util';
import TagAnnotation from './TagAnnotation';

function TagAttachAnnotationMenu({ tag, menu, onClose }) {
	const queryClient = useQueryClient();
	const newAnnotationName = useRef(null);
	const addAnnotationToTagMutation = useMutation(Client.addAnnotationToTag);
	const removeAnnotationFromTagMutation = useMutation(Client.removeAnnotationFromTag);
	const availableAnnotationsQuery = useQuery(ReactQueryUtil.availableAnnotationsKey(tag.parentId), () =>
		Client.getAvailableAnnotations(tag.parentId)
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

		addAnnotation(e, { title: newAnnotationName.current.value });
	};

	const addAnnotation = (e, annotation) => {
		e.stopPropagation();
		addAnnotationToTagMutation.mutate(
			{ tagId: tag.id, annotation: annotation },
			{
				onSuccess: (response) => {
					response.json().then((tagAnnotation) => {
						queryClient.refetchQueries({ queryKey: ReactQueryUtil.availableAnnotationsKey(tag.parentId) });
						ReactQueryUtil.updateTags(queryClient, tag.id, (currentTag) => {
							let curTagsAnnotations = currentTag.tags_annotations || [];
							curTagsAnnotations.push(tagAnnotation);

							return {
								...currentTag,
								tags_annotations: curTagsAnnotations,
							};
						});
						newAnnotationName.current.value = '';
					});
				},
			}
		);
	};

	const removeAnnotation = (e, annotation) => {
		e.stopPropagation();
		removeAnnotationFromTagMutation.mutate(
			{ tagId: tag.id, annotationId: annotation.id },
			{
				onSuccess: (response, params) => {
					queryClient.refetchQueries({ queryKey: ReactQueryUtil.availableAnnotationsKey(tag.parentId) });
					ReactQueryUtil.updateTags(queryClient, params.tagId, (currentTag) => {
						return {
							...currentTag,
							tags_annotations: currentTag.tags_annotations.filter((annotation) => {
								return annotation.id != params.annotationId;
							}),
						};
					});
				},
			}
		);
	};

	const getAnnotations = (belongToTag) => {
		if (!availableAnnotationsQuery.isSuccess) {
			return [];
		}

		return availableAnnotationsQuery.data
			.filter((annotation) => {
				if (!tag.tags_annotations) {
					return !belongToTag;
				}

				if (tag.tags_annotations.some((cur) => cur.id == annotation.id)) {
					return belongToTag;
				} else {
					return !belongToTag;
				}
			})
			.sort((a, b) => (a.title > b.title ? 1 : a.title < b.title ? -1 : 0));
	};

	return (
		<Popover
			onClose={(e, reason) => {
				if (reason == 'backdropClick' || reason == 'escapeKeyDown') {
					onClose(e);
				}
			}}
			open={true}
			anchorReference="anchorPosition"
			anchorPosition={{ top: menu.mouseY, left: menu.mouseX }}
			BackdropProps={{
				open: true,
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
					<Typography variant="body1" noWrap onClick={(e) => e.stopPropagation()}>
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
							<TagAnnotation
								key={annotation.id}
								annotation={annotation}
								selected={true}
								onRemoveClicked={removeAnnotation}
							/>
						);
					})}
					{getAnnotations(false).map((annotation) => {
						return (
							<TagAnnotation
								key={annotation.id}
								annotation={annotation}
								selected={false}
								onClick={addAnnotation}
							/>
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
						size="small"
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
