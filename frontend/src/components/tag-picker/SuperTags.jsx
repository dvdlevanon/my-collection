import { Stack } from '@mui/material';
import React from 'react';
import SuperTag from './SuperTag';

function SuperTags({ superTags, onSuperTagClicked }) {
	return (
		<Stack flexDirection="row">
			{superTags.map((tag) => {
				return <SuperTag key={tag.id} superTag={tag} onSuperTagClicked={onSuperTagClicked} />;
			})}
		</Stack>
	);
}

export default SuperTags;
