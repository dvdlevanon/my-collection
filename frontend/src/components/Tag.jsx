import AddIcon from '@mui/icons-material/Add';
import NoImageIcon from '@mui/icons-material/HideImage';
import { Box, Typography } from '@mui/material';
import React, { useState } from 'react';
import Client from '../network/client';
import TagSpeedDial from './TagSpeedDial';

function Tag({ tag, size, onTagSelected }) {
	let [optionsHidden, setOptionsHidden] = useState(true);

	const getImageUrl = () => {
		if (tag.imageUrl) {
			return Client.buildFileUrl(tag.imageUrl);
		} else {
			return Client.buildFileUrl(Client.buildInternalStoragePath('tags-image/none/1.jpg'));
		}
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
					}}
				>
					<Box component="img" src={getImageUrl()} alt={tag.title} loading="lazy" />
					{!tag.imageUrl && (
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
				{includeSpeedDial && size != 'small' && !optionsHidden && <TagSpeedDial tag={tag} />}
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
		</Box>
	);
}

export default Tag;
