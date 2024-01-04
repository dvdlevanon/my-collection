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
				{showButtons && (
					<Stack
						flexDirection="row"
						gap="5px"
						sx={{
							position: 'absolute',
							left: '0px',
							top: '0px',
							padding: '5px',
						}}
					>
						{tag && (
							<IconButton
								sx={{ border: '1px solid gray', width: '35px', height: '35px' }}
								onClick={(e) => {
									e.stopPropagation();
									e.preventDefault();
									onTagRemoved(tag);
								}}
							>
								<DeleteIcon sx={{ width: '20px', height: '20px' }} />
							</IconButton>
						)}
						<IconButton
							sx={{ border: '1px solid gray', width: '35px', height: '35px' }}
							onClick={(e) => {
								e.stopPropagation();
								e.preventDefault();
								onTagEdit(tag);
							}}
						>
							<EditIcon sx={{ width: '20px', height: '20px' }} />
						</IconButton>
					</Stack>
				)}
				<Box
					sx={{
						borderRadius: '5px',
						overflow: 'hidden',
						width: '150px',
						minWidth: '150px',
						maxHeight: '120px',
						height: '120px',
						objectFit: 'contain',
						borderRadius: '5px',
						border: '1px solid white',
						padding: '10px',
					}}
					component="img"
					src={tag == null ? TagsUtil.getNoBannerImageUrl() : TagsUtil.getTagImageUrl(tag, null, false)}
					alt={tag == null ? '' : tag.title}
					loading="lazy"
				/>
				{tag == null && (
					<NoImageIcon
						color="dark"
						sx={{
							position: 'absolute',
							width: '40%',
							height: '100%',
						}}
					/>
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
