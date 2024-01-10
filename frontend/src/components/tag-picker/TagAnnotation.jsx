import { useTheme } from '@emotion/react';
import ClearIcon from '@mui/icons-material/Clear';
import { Box, IconButton, Typography } from '@mui/material';
import React from 'react';

function TagAnnotation({ annotation, selected, onClick, onRemoveClicked }) {
	const theme = useTheme();

	return (
		<Box
			bgcolor={selected ? 'primary.dark' : 'gray'}
			color="primary.contrastText"
			sx={{
				margin: theme.spacing(1),
				padding: theme.multiSpacing(0, 1),
				display: 'flex',
				cursor: 'pointer',
			}}
			borderRadius={theme.spacing(1)}
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
					<ClearIcon sx={{ fontSize: theme.iconSize(0.8) }} />
				</IconButton>
			)}
		</Box>
	);
}

export default TagAnnotation;
