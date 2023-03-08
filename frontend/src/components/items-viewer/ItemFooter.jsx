import { Box, Tooltip, Typography } from '@mui/material';
import React from 'react';

function ItemFooter({ item }) {
	return (
		<Box
			sx={{
				display: 'flex',
				flexDirection: 'row',
				alignItems: 'center',
				padding: '10px',
				gap: '10px',
			}}
		>
			<Tooltip title={item.title} arrow followCursor>
				<Typography
					variant="caption"
					sx={{
						whiteSpace: 'nowrap',
						overflow: 'hidden',
						textOverflow: 'ellipsis',
						cursor: 'pointer',
						textAlign: 'center',
						flexGrow: 1,
					}}
				>
					{item.title}
				</Typography>
			</Tooltip>
		</Box>
	);
}

export default ItemFooter;
