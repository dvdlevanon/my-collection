import { useTheme } from '@emotion/react';
import { Chip, Link, Stack } from '@mui/material';
import React, { useState } from 'react';
import { Link as RouterLink } from 'react-router-dom';
import Client from '../../utils/client';
import TagsUtil from '../../utils/tags-util';
import ManageTagImageDialog from '../dialogs/ManageTagImageDialog';
import RemoveTagDialog from '../dialogs/RemoveTagDialog';
import TagContextMenu from '../tag-context-menu/TagContextMenu';
import TagAttachAnnotationMenu from './TagAttachAnnotationMenu';
import TagImage from './TagImage';
import TagTitle from './TagTitle';

function Tag({ tag, parent, tagDimension, selectedTit, tagLinkBuilder, onTagClicked }) {
	const [removeTagDialogOpened, setRemoveTagDialogOpened] = useState(false);
	const [attachMenuAttributes, setAttachMenuAttributes] = useState(null);
	const [manageTagImageOpened, setManageTagImageOpened] = useState(false);
	const [autoThumbnailMode, setAutoThumbnailMode] = useState(false);
	const [tagMenuProps, setTagMenuProps] = useState(null);
	const theme = useTheme();

	const onManageAttributesClicked = (e) => {
		e.stopPropagation();
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
						autoThumbnailMode={autoThumbnailMode}
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
					fontSize: theme.fontSize(1.3),
					padding: theme.multiSpacing(2, 2.5),
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
			position="relative"
			margin={parent.display_style !== 'banner' ? 'unset' : '50px'}
		>
			<Link
				onContextMenu={(e) => {
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
					onRemoveTagClicked={() => setRemoveTagDialogOpened(true)}
					withRemoveOption={true}
					withManageAttributesClicked={true}
					onManageImageClicked={() => {
						setManageTagImageOpened(true);
						setAutoThumbnailMode(false);
					}}
					onEditThumbnail={() => {
						setManageTagImageOpened(true);
						setAutoThumbnailMode(true);
					}}
					onIncludeMixClicked={() => {
						Client.includeInRandomMix(tag);
					}}
					onExcludeMixClicked={() => {
						Client.excludeFromRandomMix(tag);
					}}
				/>
			)}
		</Stack>
	);
}

export default Tag;
