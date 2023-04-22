import CancelIcon from '@mui/icons-material/Cancel';
import StopIcon from '@mui/icons-material/Stop';
import { Box, IconButton, Skeleton, Stack } from '@mui/material';
import React from 'react';

function HighlightControls({ onCancel, onDone }) {
	return (
		<Stack
			flexDirection="row"
			gap="10px"
			sx={{
				background: '#000',
				padding: '3px 10px',
				opacity: '0.7',
				borderRadius: '10px',
				position: 'absolute',
				right: 20,
				bottom: 100,
			}}
		>
			<Box sx={{ padding: '9px' }}>
				<Skeleton
					color={'red'}
					variant="circular"
					animation="pulse"
					width={25}
					height={25}
					sx={{
						backgroundColor: '#880000',
					}}
				/>
			</Box>
			<IconButton
				onClick={(e) => {
					onDone();
				}}
			>
				<StopIcon />
			</IconButton>
			<IconButton
				onClick={(e) => {
					onCancel();
				}}
			>
				<CancelIcon />
			</IconButton>
		</Stack>
	);
}

export default HighlightControls;
