import AddIcon from '@mui/icons-material/Add';

import NoImageIcon from '@mui/icons-material/HideImage';
import { Box, Typography } from '@mui/material';
import React, { useState } from 'react';
import Client from '../network/client';
import RemoveTag from './RemoveTagDialog';
import TagAttachAnnotationMenu from './TagAttachAnnotationMenu';
import TagSpeedDial from './TagSpeedDial';

function Tag({ tag, size, onTagSelected }) {
	let [optionsHidden, setOptionsHidden] = useState(true);
	let [attachMenuAttributes, setAttachMenuAttributes] = useState(null);
	let [removeTagDialogOpened, setRemoveTagDialogOpened] = useState(false);

	const getImageUrl = () => {
		if (hasImage()) {
			return Client.buildFileUrl(tag.imageUrl);
		} else {
			return Client.buildFileUrl(Client.buildInternalStoragePath('tags-image/none/1.jpg'));
		}
	};

	const hasImage = () => {
		return tag.imageUrl && tag.imageUrl != 'none';
	};

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

	const tagComponent = (placeHolderHeight, title, titleOpacity, includeSpeedDial, missingImagePlaceholder) => {
		return (
			<>
				<Box
					sx={{
						position: 'relative',
						display: 'flex',
						justifyContent: 'center',
						objectFit: 'contain',
						overflow: 'hidden',
						height: '100%',
						borderRadius: '5px',
						'&:hover': {
							filter: 'brightness(120%)',
						},
					}}
				>
					<Box
						sx={{
							borderRadius: '5px',
						}}
						component="img"
						src={getImageUrl()}
						alt={tag.title}
						loading="lazy"
					/>
					{!hasImage() && (
						<Box
							sx={{
								position: 'absolute',
								width: size == 'small' ? '50px' : '100px',
								height: placeHolderHeight,
								left: 0,
								right: 0,
								top: 0,
								bottom: 0,
								margin: 'auto',
								display: 'flex',
								flexDirection: 'column',
							}}
						>
							{missingImagePlaceholder}
						</Box>
					)}
				</Box>
				<Typography
					sx={{
						padding: '5px',
						textAlign: 'center',
						opacity: titleOpacity,
						'&:hover': {
							textDecoration: 'underline',
						},
					}}
					variant="caption"
					textAlign={'start'}
				>
					{title}
				</Typography>
				{includeSpeedDial && size != 'small' && attachMenuAttributes === null && (
					<TagSpeedDial
						hidden={optionsHidden}
						tag={tag}
						onManageAttributesClicked={onManageAttributesClicked}
						onRemoveTagClicked={() => {
							setOptionsHidden(false);
							setRemoveTagDialogOpened(true);
						}}
					/>
				)}
			</>
		);
	};

	const realTagComponent = () => {
		return tagComponent(
			size == 'small' ? '50px' : '100px',
			tag.title,
			1,
			true,
			<NoImageIcon
				color="dark"
				sx={{
					fontSize: size == 'small' ? '50px' : '100px',
				}}
			/>
		);
	};

	const newTagComponent = () => {
		return tagComponent(
			size == 'small' ? '50px' : '130px',
			'None',
			0,
			false,
			<>
				<AddIcon
					color="bright"
					sx={{
						fontSize: size == 'small' ? '50px' : '100px',
					}}
				/>
				{size != 'small' && (
					<Typography noWrap color="bright" textAlign="center" variant="button">
						New Tag
					</Typography>
				)}
			</>
		);
	};

	return (
		<Box
			sx={{
				display: 'flex',
				flexDirection: 'column',
				justifyContent: 'space-between',
				padding: '3px',
				cursor: 'pointer',
				position: 'relative',
				width: size == 'small' ? '125px' : '350px',
				height: size == 'small' ? '200px' : '500px',
			}}
			onClick={() => onTagSelected(tag)}
			onMouseEnter={() => setOptionsHidden(false)}
			onMouseLeave={() => setOptionsHidden(true)}
		>
			{tag.id > 0 ? realTagComponent() : newTagComponent()}

			{attachMenuAttributes !== null && (
				<TagAttachAnnotationMenu
					tag={tag}
					menu={attachMenuAttributes}
					onClose={() => setAttachMenuAttributes(null)}
				/>
			)}
			{removeTagDialogOpened && <RemoveTag tag={tag} onClose={() => setRemoveTagDialogOpened(false)} />}
		</Box>
	);
}

export default Tag;
