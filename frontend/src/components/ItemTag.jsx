import CloseIcon from '@mui/icons-material/Close';
import { Box, IconButton } from '@mui/material';
import styles from './ItemTag.module.css';

function ItemTag({ tag, onRemoveClicked }) {
	const removeHandler = (e) => {
		e.stopPropagation();
		onRemoveClicked(tag);
	};

	return (
		<Box
			sx={{
				backgroundColor: 'primary.main',
				'&:hover': {
					backgroundColor: 'primary.main',
					opacity: [0.9, 0.8, 0.7],
				},
			}}
			className={styles.item_tag}
		>
			<IconButton onClick={(e) => removeHandler(e)}>
				<CloseIcon />
			</IconButton>
			{tag.title}
		</Box>
	);
}

export default ItemTag;
