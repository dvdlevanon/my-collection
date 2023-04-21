import { Typography } from '@mui/material';
import React from 'react';

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
			{tag.title}
		</Typography>
	);
}

export default TagTitle;
