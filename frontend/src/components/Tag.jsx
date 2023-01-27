import { Box, Typography } from '@mui/material';
import { useState } from 'react';
import Client from '../network/client';
import styles from './Tag.module.css';
import TagSpeedDial from './TagSpeedDial';

function Tag({ tag, size, onTagSelected }) {
	let [optionsHidden, setOptionsHidden] = useState(true);

	const getTagClasses = () => {
		if (size != 'small') {
			return styles.tag + ' ' + styles.big_tag;
		} else {
			return styles.tag;
		}
	};

	const getImageUrl = () => {
		if (tag.imageUrl) {
			return Client.buildFileUrl(tag.imageUrl);
		} else {
			return 'empty';
		}
	};

	return (
		<Box
			className={getTagClasses()}
			onClick={() => onTagSelected(tag)}
			onMouseEnter={() => setOptionsHidden(false)}
			onMouseLeave={() => setOptionsHidden(true)}
		>
			{size != 'small' && <img className={styles.image} src={getImageUrl()} alt={tag.title} loading="lazy" />}
			<Typography
				className={styles.title}
				sx={{
					'&:hover': {
						textDecoration: 'underline',
					},
				}}
				variant="caption"
				textAlign={'start'}
			>
				{tag.title}
			</Typography>
			{size != 'small' && !optionsHidden && <TagSpeedDial tag={tag} />}
		</Box>
	);
}

export default Tag;
