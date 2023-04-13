import AddIcon from '@mui/icons-material/Add';
import CloseIcon from '@mui/icons-material/Close';
import { Dialog, DialogContent, DialogTitle, IconButton, Stack, TextField, Tooltip } from '@mui/material';
import React, { useRef, useState } from 'react';
import { useQueryClient } from 'react-query';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import CategoriesChooser from '../categories-chooser/CategoriesChooser';

function AddTagDialog({ parentId, verb, onTagAdded, onClose, shouldSelectParent, initialText }) {
	const unknownCategoryId = -1;
	const queryClient = useQueryClient();
	const newTagName = useRef(null);
	const [selectedCategory, setSelectedCategory] = useState(-1);

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
			open={true}
		>
			<DialogTitle variant="h6">
				Add {verb}
				<IconButton
					sx={{
						position: 'absolute',
						top: '0px',
						right: '0px',
						margin: '10px',
					}}
					onClick={() => onClose()}
				>
					<CloseIcon />
				</IconButton>
			</DialogTitle>
			<DialogContent
				sx={{
					display: 'flex',
					flexDirection: 'column',
					gap: '10px',
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
				<Stack flexDirection="row" gap="10px">
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
							<AddIcon />
						</IconButton>
					</Tooltip>
				</Stack>
			</DialogContent>
		</Dialog>
	);
}

export default AddTagDialog;
