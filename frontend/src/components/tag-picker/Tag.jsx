import AddIcon from '@mui/icons-material/Add';

import NoImageIcon from '@mui/icons-material/HideImage';
import { Box, Link, Typography } from '@mui/material';
import React, { useState } from 'react';
import { Link as RouterLink } from 'react-router-dom';
import Client from '../../utils/client';
import GalleryUrlParams from '../../utils/gallery-url-params';
import TagsUtil from '../../utils/tags-util';
import ManageTagImageDialog from '../dialogs/ManageTagImageDialog';
import RemoveTagDialog from '../dialogs/RemoveTagDialog';
import TagAttachAnnotationMenu from './TagAttachAnnotationMenu';
import TagSpeedDial from './TagSpeedDial';
function Tag({ tag, size, selectedTit, onTagSelected }) {
	let [optionsHidden, setOptionsHidden] = useState(true);
	let [attachMenuAttributes, setAttachMenuAttributes] = useState(null);
	let [removeTagDialogOpened, setRemoveTagDialogOpened] = useState(false);
	let [manageTagImageOpened, setManageTagImageOpened] = useState(false);

	const getImageUrl = () => {
		if (TagsUtil.isDirectoriesCategory(tag.parentId)) {
			return Client.buildFileUrl(Client.buildInternalStoragePath('tags-image/directory/directory.png'));
		} else if (TagsUtil.isDailymixCategory(tag.parentId)) {
			return Client.buildFileUrl(Client.buildInternalStoragePath('tags-image/dailymix/dailymix.png'));
		} else {
			return getRegularTagImageUrl();
		}
	};

	const getRegularTagImageUrl = () => {
		if (selectedTit && tag.images) {
			let selectedImage = tag.images.find((image) => image.imageType == selectedTit.id);
			if (selectedImage && selectedImage.url) {
				return Client.buildFileUrl(selectedImage.url);
			}
		}

		if (tag.images) {
			for (let i = 0; i < tag.images.length; i++) {
				if (tag.images[i].url) {
					return Client.buildFileUrl(tag.images[i].url);
				}
			}
		}

		return Client.buildFileUrl(Client.buildInternalStoragePath('tags-image/none/1.jpg'));
	};

	const hasImage = () => {
		return tag.images && tag.images.length > 0;
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
							borderColor: 'primary.main',
						},
					}}
				>
					<Box
						sx={{
							borderRadius: '5px',
							objectFit: 'contain',
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
								width: size == 'small' ? '70px' : '100px',
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
					noWrap
					variant="caption"
					textAlign={'start'}
				>
					{title}
				</Typography>
				{includeSpeedDial && attachMenuAttributes === null && (
					<TagSpeedDial
						hidden={optionsHidden}
						tag={tag}
						onManageImageClicked={() => setManageTagImageOpened(true)}
						onManageAttributesClicked={onManageAttributesClicked}
						onRemoveTagClicked={() => setRemoveTagDialogOpened(true)}
					/>
				)}
			</>
		);
	};

	const realTagComponent = () => {
		return tagComponent(
			size == 'small' ? '70px' : '100px',
			tag.title,
			1,
			true,
			<NoImageIcon
				color="dark"
				sx={{
					fontSize: size == 'small' ? '70px' : '100px',
				}}
			/>
		);
	};

	const newTagComponent = () => {
		return tagComponent(
			size == 'small' ? '70px' : '130px',
			'None',
			0,
			false,
			<>
				<AddIcon
					color="bright"
					sx={{
						fontSize: size == 'small' ? '70px' : '100px',
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

	const specialTagComponent = () => {
		return tagComponent('50px', tag.title, 1, false, <></>);
	};

	const getTagComponent = () => {
		if (TagsUtil.isSpecialCategory(tag.parentId)) {
			return specialTagComponent();
		} else if (tag.id > 0) {
			return realTagComponent();
		} else {
			return newTagComponent();
		}
	};

	const getWidth = () => {
		if (TagsUtil.isSpecialCategory(tag.parentId)) {
			return '200px';
		} else if (size == 'small') {
			return '225px';
		} else {
			return '350px';
		}
	};

	const getHeight = () => {
		if (TagsUtil.isSpecialCategory(tag.parentId)) {
			return '200px';
		} else if (size == 'small') {
			return '300px';
		} else {
			return '500px';
		}
	};

	return (
		<Link component={RouterLink} to={'/?' + GalleryUrlParams.buildUrlParams(tag.id)}>
			<Box
				sx={{
					display: 'flex',
					flexDirection: 'column',
					justifyContent: 'space-between',
					padding: '3px',
					cursor: 'pointer',
					position: 'relative',
					width: getWidth(),
					height: getHeight(),
				}}
				onClick={(e) => {
					e.preventDefault();
					onTagSelected(tag);
				}}
				onMouseEnter={() => setOptionsHidden(false)}
				onMouseLeave={() => setOptionsHidden(true)}
			>
				{getTagComponent()}

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
			</Box>
		</Link>
	);
}

export default Tag;
