import { Box, Tooltip, Typography } from '@mui/material';
import React from 'react';

function ItemFooter({ item }) {
	const getFormattedDuration = () => {
		if (!item.duration_seconds) {
			return '00:00';
		}

		if (item.duration_seconds < 60 * 60) {
			return new Date(item.duration_seconds * 1000).toISOString().slice(14, 19);
		} else {
			return new Date(item.duration_seconds * 1000).toISOString().slice(11, 19);
		}
	};

	return (
		<Box
			sx={{
				display: 'flex',
				flexDirection: 'row',
				gap: '10px',
				alignItems: 'center',
				justifyContent: 'center',
				height: '50px',
			}}
		>
			<Typography
				variant="caption"
				sx={{
					padding: '0px 3px',
					borderWidth: '1px',
					borderColor: 'bright.main',
					borderStyle: 'solid',
					borderRadius: '3px',
					color: 'bright.main',
					verticalAlign: 'middle',
					margin: '10px',
				}}
			>
				{getFormattedDuration()}
			</Typography>
			<Tooltip title={item.title} arrow followCursor>
				<Typography
					variant="caption"
					sx={{
						whiteSpace: 'nowrap',
						overflow: 'hidden',
						textOverflow: 'ellipsis',
						cursor: 'pointer',
						maxWidth: '450px',
						textAlign: 'center',
						padding: '5px',
						color: 'primary.light',
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
