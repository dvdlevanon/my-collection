import ClearIcon from '@mui/icons-material/Clear';
import { Box, IconButton, Typography } from '@mui/material';
import React from 'react';

function TagAnnotation({ annotation, selected, onClick, onRemoveClicked }) {
	return (
		<Box
			bgcolor={selected ? 'primary.dark' : 'gray'}
			color="primary.contrastText"
			sx={{
				margin: '10px',
				padding: '0px 10px',
				display: 'flex',
				cursor: 'pointer',
			}}
			borderRadius="10px"
			onClick={(e) => {
				if (onClick) {
					onClick(e, annotation);
				}
			}}
			key={annotation.id}
		>
			<Box sx={{ display: 'flex', alignItems: 'center' }}>
				<Typography noWrap sx={{ flexGrow: 1 }} variant="body2">
					{annotation.title}
				</Typography>
			</Box>
			{onRemoveClicked && (
				<IconButton onClick={(e) => onRemoveClicked(e, annotation)} size="small">
					<ClearIcon sx={{ fontSize: '15px' }} />
				</IconButton>
			)}
		</Box>
	);
}

export default TagAnnotation;
