import { useTheme } from '@emotion/react';
import AddIcon from '@mui/icons-material/Add';
import CloseIcon from '@mui/icons-material/Close';
import { Divider, IconButton, Popover, Stack, TextField, Typography } from '@mui/material';
import { Box } from '@mui/system';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useRef } from 'react';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import TagAnnotation from './TagAnnotation';

function TagAttachAnnotationMenu({ tag, menu, onClose }) {
	const queryClient = useQueryClient();
	const newAnnotationName = useRef(null);
	const addAnnotationToTagMutation = useMutation(Client.addAnnotationToTag);
	const removeAnnotationFromTagMutation = useMutation(Client.removeAnnotationFromTag);
	const theme = useTheme();
	const availableAnnotationsQuery = useQuery({
		queryKey: ReactQueryUtil.availableAnnotationsKey(tag.parentId),

		queryFn: () => Client.getAvailableAnnotations(tag.parentId),
	});

	const handleClose = (e) => {
		e.preventDefault();
		e.stopPropagation();
		onClose();
	};

	const addNewAnnotation = (e) => {
		e.preventDefault();
		e.stopPropagation();
		if (newAnnotationName.current.value == '') {
			return;
		}

		addAnnotation(e, { title: newAnnotationName.current.value });
	};

	const addAnnotation = (e, annotation) => {
		e.preventDefault();
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
		e.preventDefault();
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
					gap: theme.spacing(1),
					flexDirection: 'column',
					padding: theme.spacing(1),
				},
				onClick: (e) => {
					e.preventDefault();
					e.stopPropagation();
				},
			}}
		>
			<Box>
				<Box
					onClick={(e) => {
						e.preventDefault();
						e.stopPropagation();
					}}
					sx={{
						display: 'flex',
						gap: theme.spacing(1),
						alignItems: 'center',
					}}
				>
					<IconButton onClick={handleClose}>
						<CloseIcon sx={{ fontSize: theme.iconSize(1) }} />
					</IconButton>
					<Typography
						variant="body1"
						noWrap
						onClick={(e) => {
							e.preventDefault();
							e.stopPropagation();
						}}
					>
						{tag.title} Annotations
					</Typography>
				</Box>
				<Divider />
				<Stack flexDirection="row" flexWrap="wrap">
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
				</Stack>
				<Box
					sx={{
						display: 'flex',
						gap: theme.spacing(1),
						justifyContent: 'center',
						alignItems: 'center',
					}}
					onClick={(e) => {
						e.preventDefault();
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
						<AddIcon sx={{ fontSize: theme.iconSize(1) }} />
					</IconButton>
				</Box>
			</Box>
		</Popover>
	);
}

export default TagAttachAnnotationMenu;
