import NoImageIcon from '@mui/icons-material/HideImage';
import { Box } from '@mui/material';
import React from 'react';
import TagsUtil from '../../utils/tags-util';

function TagImage({ tag, selectedTit, onClick, imgSx }) {
	return (
		<Box
			onClick={onClick}
			sx={{
				position: 'relative',
				display: 'flex',
				overflow: 'hidden',
				objectFit: 'contain',
				justifyContent: 'center',
				height: '100%',
				width: '100%',
				cursor: 'pointer',
				borderRadius: '10px',
				'&:hover': {
					filter: 'brightness(120%)',
				},
			}}
		>
			<Box
				sx={{
					borderRadius: '5px',
					overflow: 'hidden',
					...imgSx,
				}}
				component="img"
				src={TagsUtil.getTagImageUrl(tag, selectedTit)}
				alt={tag.title}
				loading="lazy"
			/>
			{!TagsUtil.hasImage(tag) && (
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
}

export default TagImage;
