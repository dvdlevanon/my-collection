import { Menu, MenuItem, Typography } from '@mui/material';
import React, { useState } from 'react';
import { useQuery } from 'react-query';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import TagsUtil from '../../utils/tags-util';
import AddTagDialog from '../dialogs/AddTagDialog';

function ItemTitle({ item, variant, onTagAdded }) {
	const tagsQuery = useQuery(ReactQueryUtil.TAGS_KEY, Client.getTags);
	const [menuAchrosEl, setMenuAchrosEl] = useState(null);
	const [menuX, setMenuX] = useState(null);
	const [menuY, setMenuY] = useState(null);
	const [selectedText, setSelectedText] = useState(null);
	const [addTagDialogOpened, setAddTagDialogOpened] = useState(false);
	const [menuDelayTimer, setMenuDelayTimer] = useState(0);

	const getTagByTitle = (tagTitle) => {
		let tags = tagsQuery.data.filter((tag) => {
			return tagTitle.toLowerCase() == tag.title.toLowerCase() && !TagsUtil.isSpecialCategory(tag.parentId);
		});

		if (tags.length == 1) {
			return tags[0];
		}

		return null;
	};

	const getTagMenuText = () => {
		let title = TagsUtil.normalizeTagTitle(selectedText);

		if (getTagByTitle(title)) {
			return "Add to '" + title + "'";
		} else {
			return "Open create tag dialog for '" + title + "'";
		}
	};

	const tagMenuClicked = () => {
		let title = TagsUtil.normalizeTagTitle(selectedText);
		let existingTag = getTagByTitle(title);

		if (getTagByTitle(title)) {
			if (onTagAdded) {
				onTagAdded(existingTag);
			} else {
				alert('onTagAdded is undefined');
			}
			setMenuAchrosEl(null);
		} else {
			setAddTagDialogOpened(true);
			setMenuAchrosEl(null);
		}
	};

	return (
		<>
			<Typography
				variant={variant}
				onMouseUp={(e) => {
					let event = e;
					console.log('up');
					setMenuDelayTimer(
						setTimeout(() => {
							let selection = window.getSelection();
							let selectedText = selection.toString();

							if (!selectedText) {
								setSelectedText(null);
								setMenuX(event.clientX);
								setMenuY(event.clientY);
							} else {
								setSelectedText(selectedText);
								let selectionRect = selection.getRangeAt(0).getBoundingClientRect();
								setMenuX(selectionRect.left);
								setMenuY(selectionRect.bottom);
							}

							setMenuAchrosEl(event.target);
							setMenuDelayTimer(0);
						}, 200)
					);
				}}
				onMouseDown={() => {
					if (menuDelayTimer != 0) {
						clearTimeout(menuDelayTimer);
						setMenuDelayTimer(0);
					}
				}}
			>
				{item.title}
			</Typography>
			{menuAchrosEl != null && tagsQuery.isSuccess && (
				<Menu
					open={menuAchrosEl != null}
					anchorEl={menuAchrosEl}
					onClose={() => setMenuAchrosEl(null)}
					anchorPosition={{ left: menuX, top: menuY }}
					anchorReference="anchorPosition"
				>
					{selectedText && <MenuItem onClick={tagMenuClicked}>{getTagMenuText()}</MenuItem>}
					{selectedText && (
						<MenuItem
							onClick={(e) => {
								navigator.clipboard.writeText(TagsUtil.normalizeTagTitle(selectedText));
								setMenuAchrosEl(null);
							}}
						>
							Copy '{TagsUtil.normalizeTagTitle(selectedText)}'
						</MenuItem>
					)}
					<MenuItem
						onClick={(e) => {
							Client.getItemLocation(item.id).then((location) => {
								navigator.clipboard.writeText(`"${location.url}"`);
							});

							setMenuAchrosEl(null);
						}}
					>
						Copy file location
					</MenuItem>
				</Menu>
			)}
			{addTagDialogOpened && (
				<AddTagDialog
					parentId={null}
					verb="Tag"
					initialText={TagsUtil.normalizeTagTitle(selectedText)}
					shouldSelectParent={true}
					onClose={() => setAddTagDialogOpened(false)}
					onTagAdded={(newTag) => {
						if (onTagAdded) {
							onTagAdded(newTag);
						}
					}}
				/>
			)}
		</>
	);
}

export default ItemTitle;
