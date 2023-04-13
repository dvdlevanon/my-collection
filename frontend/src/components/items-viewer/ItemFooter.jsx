import { Box, Tooltip } from '@mui/material';
import React from 'react';
import ItemTitle from './ItemTitle';

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
				<Box
					sx={{
						whiteSpace: 'nowrap',
						overflow: 'hidden',
						textOverflow: 'ellipsis',
						textAlign: 'center',
						flexGrow: 1,
					}}
				>
					<ItemTitle item={item} variant="caption" />
				</Box>
			</Tooltip>
		</Box>
	);
}

export default ItemFooter;
