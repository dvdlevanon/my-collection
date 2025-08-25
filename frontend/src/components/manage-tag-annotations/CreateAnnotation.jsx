import AddIcon from '@mui/icons-material/Add';
import { Box, IconButton, TextField, useTheme } from '@mui/material';
import { useRef } from 'react';
import { useAnnotations } from './useAnnotations';

function CreateAnnotation() {
	const theme = useTheme();
	const newAnnotationName = useRef(null);
	const { addAnnotation } = useAnnotations();

	const addNewAnnotation = (e) => {
		e.preventDefault();
		e.stopPropagation();
		if (newAnnotationName.current.value == '') {
			return;
		}

		addAnnotation({ title: newAnnotationName.current.value }, () => (newAnnotationName.current.value = ''));
	};

	return (
		<Box
			sx={{
				display: 'flex',
				gap: theme.spacing(1),
				justifyContent: 'center',
				alignItems: 'center',
			}}
			onClick={(e) => {
				e.preventDefault();
				e.stopPropagation();
			}}
		>
			<TextField
				onKeyDown={(e) => {
					if (e.key == 'Enter') {
						addNewAnnotation(e);
					}
				}}
				size="small"
				autoFocus
				inputRef={newAnnotationName}
				placeholder="New Annotation Name..."
				sx={{
					flexGrow: 1,
				}}
			></TextField>
			<IconButton onClick={(e) => addNewAnnotation(e)} sx={{ alignSelf: 'center' }}>
				<AddIcon sx={{ fontSize: theme.iconSize(1) }} />
			</IconButton>
		</Box>
	);
}

export default CreateAnnotation;
