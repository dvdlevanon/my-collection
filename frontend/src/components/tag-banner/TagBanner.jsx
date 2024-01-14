import { useTheme } from '@emotion/react';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';
import NoImageIcon from '@mui/icons-material/HideImage';
import { Box, IconButton, Stack } from '@mui/material';
import React, { useState } from 'react';
import { Link, Link as RouterLink } from 'react-router-dom';
import GalleryUrlParams from '../../utils/gallery-url-params';
import TagsUtil from '../../utils/tags-util';

function TagBanner({ tag, onTagRemoved, onTagEdit }) {
	const [showButtons, setShowButtons] = useState(false);
	const theme = useTheme();

	const removeUrlParams = (url) => {
		const indexOfQuestionMark = url.indexOf('?');

		if (indexOfQuestionMark === -1) {
			return url;
		}

		return url.substring(0, indexOfQuestionMark);
	};

	const getBannerComponent = () => {
		return (
			<Box
				sx={{
					position: 'relative',
					display: 'flex',
					overflow: 'hidden',
					objectFit: 'contain',
					justifyContent: 'center',
					cursor: 'pointer',
				}}
				onMouseEnter={() => setShowButtons(true)}
				onMouseLeave={() => setShowButtons(false)}
			>
				<Box
					sx={{
						overflow: 'hidden',
						width: theme.iconSize(6),
						minWidth: theme.iconSize(6),
						maxHeight: theme.iconSize(5),
						height: theme.iconSize(5),
						objectFit: tag == null ? 'auto' : 'contain',
						borderRadius: theme.spacing(0.5),
						border: theme.border(1, 'solid', theme.palette.text.primary),
						padding: tag == null ? '0px' : theme.spacing(1),
					}}
					component="img"
					src={tag == null ? TagsUtil.getNoBannerImageUrl() : TagsUtil.getTagImageUrl(tag, null, false)}
					alt={tag == null ? '' : tag.title}
					loading="lazy"
				/>
				{tag == null && (
					<NoImageIcon
						sx={{
							color: theme.palette.primary.light,
							position: 'absolute',
							width: '40%',
							height: '100%',
							fontSize: theme.iconSize(1),
						}}
					/>
				)}
				{showButtons && (
					<Stack
						flexDirection="row"
						gap={theme.spacing(0.5)}
						sx={{
							position: 'absolute',
							left: '0',
							top: '0',
							padding: theme.spacing(0.5),
						}}
					>
						{tag && (
							<IconButton
								sx={{
									border: theme.border(1, 'solid', theme.palette.text.primary),
									width: theme.iconSize(1.4),
									height: theme.iconSize(1.4),
								}}
								onClick={(e) => {
									e.stopPropagation();
									e.preventDefault();
									onTagRemoved(tag);
								}}
							>
								<DeleteIcon sx={{ width: theme.iconSize(0.9), height: theme.iconSize(0.9) }} />
							</IconButton>
						)}
						<IconButton
							color="secondary"
							sx={{
								border: theme.border(1, 'solid', theme.palette.secondary.main),
								width: theme.iconSize(1.4),
								height: theme.iconSize(1.4),
							}}
							onClick={(e) => {
								e.stopPropagation();
								e.preventDefault();
								onTagEdit(tag);
							}}
						>
							<EditIcon sx={{ width: theme.iconSize(0.9), height: theme.iconSize(0.9) }} />
						</IconButton>
					</Stack>
				)}
			</Box>
		);
	};

	return (
		<>
			{(!tag && getBannerComponent()) || (
				<Link target="_blank" component={RouterLink} to={'/?' + GalleryUrlParams.buildUrlParams(tag.id)}>
					{getBannerComponent()}
				</Link>
			)}
		</>
	);
}

export default TagBanner;
