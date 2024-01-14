import { useTheme } from '@emotion/react';
import NoImageIcon from '@mui/icons-material/HideImage';
import { Box } from '@mui/material';
import React from 'react';
import TagsUtil from '../../utils/tags-util';

function TagImage({ tag, selectedTit, onClick, imgSx }) {
	const theme = useTheme();

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
				borderRadius: theme.spacing(1),
				'&:hover': {
					filter: 'brightness(105%)',
				},
			}}
		>
			<Box
				sx={{
					borderRadius: theme.spacing(0.5),
					overflow: 'hidden',
					...imgSx,
				}}
				component="img"
				src={TagsUtil.getTagImageUrl(tag, selectedTit, false)}
				alt={tag.title}
				loading="lazy"
			/>
			{!TagsUtil.hasImage(tag) && (
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
		</Box>
	);
}

export default TagImage;
