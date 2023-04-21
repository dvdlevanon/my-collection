import { Link, Stack } from '@mui/material';
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
			// maxWidth={tagDimension.width}
			// maxHeight={tagDimension.height}
			width={tagDimension.width}
			height={tagDimension.height}
			onMouseEnter={() => setOptionsHidden(false)}
			onMouseLeave={() => setOptionsHidden(true)}
			position="relative"
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
				<TagImage
					tag={tag}
					selectedTit={selectedTit}
					onClick={(e) => {
						e.preventDefault();
						onTagClicked(tag);
					}}
				/>
			</Link>
			<TagTitle tag={tag} />
			{optionsComponents()}
		</Stack>
	);
}

export default Tag;
