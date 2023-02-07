import CloseIcon from '@mui/icons-material/Close';
import { Box, IconButton } from '@mui/material';

function ItemTag({ tag, onRemoveClicked }) {
	const removeHandler = (e) => {
		e.stopPropagation();
		onRemoveClicked(tag);
	};

	return (
		<Box
			sx={{
				borderRadius: '7px',
				cursor: 'pointer',
				padding: '3px',
				backgroundColor: 'primary.main',
				'&:hover': {
					backgroundColor: 'primary.main',
					opacity: [0.9, 0.8, 0.7],
				},
			}}
		>
			<IconButton onClick={(e) => removeHandler(e)}>
				<CloseIcon />
			</IconButton>
			{tag.title}
		</Box>
	);
}

export default ItemTag;
