import { Typography } from '@mui/material';
import React from 'react';
import TagsUtil from '../../utils/tags-util';

function TagTitle({ tag }) {
	return (
		<Typography
			sx={{
				padding: '10px',
				textAlign: 'center',
			}}
			noWrap
			variant="caption"
			textAlign={'start'}
		>
			{tag.title + ' (' + TagsUtil.itemsCount(tag) + ')'}
		</Typography>
	);
}

export default TagTitle;
