import { Chip, Link, Stack } from '@mui/material';
import React, { useState } from 'react';
import { Link as RouterLink } from 'react-router-dom';
import TagsUtil from '../../utils/tags-util';
import ManageTagImageDialog from '../dialogs/ManageTagImageDialog';
import RemoveTagDialog from '../dialogs/RemoveTagDialog';
import TagContextMenu from '../tag-context-menu/TagContextMenu';
import TagAttachAnnotationMenu from './TagAttachAnnotationMenu';
import TagImage from './TagImage';
import TagTitle from './TagTitle';

function Tag({ tag, parent, tagDimension, selectedTit, tagLinkBuilder, onTagClicked }) {
	const [optionsHidden, setOptionsHidden] = useState(true);
	const [attachMenuAttributes, setAttachMenuAttributes] = useState(null);
	const [removeTagDialogOpened, setRemoveTagDialogOpened] = useState(false);
	const [manageTagImageOpened, setManageTagImageOpened] = useState(false);
	const [tagMenuProps, setTagMenuProps] = useState(null);

	const onManageAttributesClicked = (e) => {
		e.stopPropagation();
		setOptionsHidden(false);
		setAttachMenuAttributes(
			attachMenuAttributes === null
				? {
						mouseX: e.clientX + 2,
						mouseY: e.clientY - 6,
				  }
				: null
		);
	};

	const optionsComponents = () => {
		return (
			<>
				{attachMenuAttributes !== null && (
					<TagAttachAnnotationMenu
						tag={tag}
						menu={attachMenuAttributes}
						onClose={() => setAttachMenuAttributes(null)}
					/>
				)}
				{removeTagDialogOpened && <RemoveTagDialog tag={tag} onClose={() => setRemoveTagDialogOpened(false)} />}
				{manageTagImageOpened && (
					<ManageTagImageDialog
						tag={tag}
						autoThumbnailMode={false}
						onClose={() => setManageTagImageOpened(false)}
					/>
				)}
			</>
		);
	};

	const tagChipComponent = () => {
		return (
			<Chip
				label={tag.title + ' (' + TagsUtil.itemsCount(tag) + ')'}
				variant="outlined"
				onClick={(e) => {
					e.preventDefault();
					onTagClicked(tag);
				}}
				sx={{
					cursor: 'pointer',
					fontSize: '20px',
					padding: '20px 25px',
				}}
			></Chip>
		);
	};

	const tagImageComponent = () => {
		return (
			<TagImage
				tag={tag}
				selectedTit={selectedTit}
				onClick={(e) => {
					e.preventDefault();
					onTagClicked(tag);
				}}
				imgSx={{
					overflow: parent.display_style !== 'banner' ? 'visible' : 'hidden',
					objectFit: parent.display_style !== 'banner' ? 'auto' : 'contain',
				}}
			/>
		);
	};

	return (
		<Stack
			maxWidth={tagDimension.width}
			maxHeight={tagDimension.height}
			width={parent.display_style !== 'banner' ? tagDimension.width : 'auto'}
			height={parent.display_style !== 'banner' ? tagDimension.height : 'auto'}
			onMouseEnter={() => setOptionsHidden(false)}
			onMouseLeave={() => setOptionsHidden(true)}
			position="relative"
			margin={parent.display_style !== 'banner' ? 'unset' : '50px'}
		>
			<Link
				onContextMenu={(e) => {
					e.preventDefault();
					e.preventDefault();
					setTagMenuProps({
						anchor: e.target,
						left: e.clientX,
						top: e.clientY,
					});
				}}
				component={RouterLink}
				to={tagLinkBuilder(tag)}
				sx={{
					height: '100%',
					width: '100%',
					overflow: 'hidden',
				}}
			>
				{(parent.display_style === 'chip' && tagChipComponent()) || tagImageComponent()}
			</Link>
			{parent.display_style !== 'banner' && parent.display_style !== 'chip' && <TagTitle tag={tag} />}
			{optionsComponents()}
			{tagMenuProps != null && (
				<TagContextMenu
					tag={tag}
					menuAnchorEl={tagMenuProps.anchor}
					menuPosition={{ top: tagMenuProps.top, left: tagMenuProps.left }}
					onClose={() => setTagMenuProps(null)}
					onManageAttributesClicked={onManageAttributesClicked}
					onManageImageClicked={setManageTagImageOpened}
					onRemoveTagClicked={() => setRemoveTagDialogOpened(true)}
				/>
			)}
		</Stack>
	);
}

export default Tag;
