import { Chip, Link, Stack } from '@mui/material';
import React, { useState } from 'react';
import { Link as RouterLink } from 'react-router-dom';
import ManageTagImageDialog from '../dialogs/ManageTagImageDialog';
import RemoveTagDialog from '../dialogs/RemoveTagDialog';
import TagAttachAnnotationMenu from './TagAttachAnnotationMenu';
import TagImage from './TagImage';
import TagSpeedDial from './TagSpeedDial';
import TagTitle from './TagTitle';

function Tag({ tag, parent, tagDimension, selectedTit, tagLinkBuilder, onTagClicked }) {
	const [optionsHidden, setOptionsHidden] = useState(true);
	const [attachMenuAttributes, setAttachMenuAttributes] = useState(null);
	const [removeTagDialogOpened, setRemoveTagDialogOpened] = useState(false);
	const [manageTagImageOpened, setManageTagImageOpened] = useState(false);

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
				{attachMenuAttributes === null && (
					<TagSpeedDial
						hidden={optionsHidden}
						tag={tag}
						onManageImageClicked={() => setManageTagImageOpened(true)}
						onManageAttributesClicked={onManageAttributesClicked}
						onRemoveTagClicked={() => setRemoveTagDialogOpened(true)}
					/>
				)}
				{attachMenuAttributes !== null && (
					<TagAttachAnnotationMenu
						tag={tag}
						menu={attachMenuAttributes}
						onClose={() => setAttachMenuAttributes(null)}
					/>
				)}
				{removeTagDialogOpened && <RemoveTagDialog tag={tag} onClose={() => setRemoveTagDialogOpened(false)} />}
				{manageTagImageOpened && (
					<ManageTagImageDialog tag={tag} onClose={() => setManageTagImageOpened(false)} />
				)}
			</>
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
				component={RouterLink}
				to={tagLinkBuilder(tag)}
				sx={{
					height: '100%',
					width: '100%',
					overflow: 'hidden',
				}}
			>
				{(parent.display_style === 'chip' && (
					<Chip
						label={tag.title}
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
					/>
				)) || (
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
				)}
			</Link>
			{parent.display_style !== 'banner' && parent.display_style !== 'chip' && <TagTitle tag={tag} />}
			{parent.display_style !== 'chip' && optionsComponents()}
		</Stack>
	);
}

export default Tag;
