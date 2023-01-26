import { Box } from '@mui/material';
import React from 'react';
import SuperTag from './SuperTag';
import styles from './SuperTags.module.css';

function SuperTags({ superTags, onSuperTagClicked }) {
	return (
		<Box className={styles.super_tags}>
			{superTags.map((tag) => {
				return <SuperTag key={tag.id} superTag={tag} onSuperTagClicked={onSuperTagClicked} />;
			})}
		</Box>
	);
}

export default SuperTags;
