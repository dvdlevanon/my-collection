import { useTheme } from '@emotion/react';
import AddIcon from '@mui/icons-material/Add';
import CloseIcon from '@mui/icons-material/Close';
import { Dialog, DialogContent, DialogTitle, IconButton, Stack, TextField, Tooltip } from '@mui/material';
import React, { useRef, useState } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import CategoriesChooser from '../categories-chooser/CategoriesChooser';

function AddTagDialog({ open, parentId, verb, onTagAdded, onClose, shouldSelectParent, initialText }) {
	const unknownCategoryId = -1;
	const queryClient = useQueryClient();
	const newTagName = useRef(null);
	const [selectedCategory, setSelectedCategory] = useState(-1);
	const theme = useTheme();

	const addTag = (e) => {
		if (!newTagName.current.value) {
			return;
		}

		if (shouldSelectParent && selectedCategory < 1) {
			return;
		}

		let newTag = {
			title: newTagName.current.value,
		};

		if (parentId != null) {
			newTag.parentId = parentId;
		} else if (shouldSelectParent) {
			newTag.parentId = selectedCategory;
		}

		Client.createTag(newTag)
			.then((response) => response.json())
			.then((newTag) => {
				queryClient.refetchQueries({ queryKey: ReactQueryUtil.TAGS_KEY });
				onClose();
				if (onTagAdded) {
					onTagAdded(newTag);
				}
			});
	};

	return (
		<Dialog
			onClose={(e, reason) => {
				if (reason == 'backdropClick' || reason == 'escapeKeyDown') {
					onClose();
				}
			}}
			open={open}
		>
			<DialogTitle variant="h6">
				Add {verb}
				<IconButton
					sx={{
						position: 'absolute',
						top: '0',
						right: '0',
						margin: theme.spacing(1),
					}}
					onClick={() => onClose()}
				>
					<CloseIcon sx={{ fontSize: theme.iconSize(1) }} />
				</IconButton>
			</DialogTitle>
			<DialogContent
				sx={{
					display: 'flex',
					flexDirection: 'column',
					gap: theme.spacing(1),
				}}
			>
				{shouldSelectParent && (
					<CategoriesChooser
						multiselect={false}
						allowToCreate={false}
						placeholder="Select category"
						selectedIds={selectedCategory == null ? [unknownCategoryId] : [selectedCategory]}
						setCategories={(categoryId) => {
							setSelectedCategory(categoryId);
						}}
					/>
				)}
				<Stack flexDirection="row" gap={theme.spacing(1)}>
					<TextField
						autoFocus
						onKeyDown={(e) => {
							if (e.key == 'Enter') {
								addTag(e);
							}
						}}
						size="small"
						placeholder={verb + ' Name'}
						inputRef={newTagName}
						defaultValue={initialText}
					></TextField>
					<Tooltip title="Add">
						<IconButton onClick={(e) => addTag(e)} sx={{ alignSelf: 'center' }}>
							<AddIcon sx={{ fontSize: theme.iconSize(1) }} />
						</IconButton>
					</Tooltip>
				</Stack>
			</DialogContent>
		</Dialog>
	);
}

export default AddTagDialog;
