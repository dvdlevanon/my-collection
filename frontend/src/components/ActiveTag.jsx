import CloseIcon from '@mui/icons-material/Close';
import { Box, IconButton, Typography } from '@mui/material';

function ActiveTag({ tag, onTagDeactivated, onTagSelected, onTagDeselected }) {
	const onTagClicked = (e) => {
		if (tag.selected) {
			onTagDeselected(tag);
		} else {
			onTagSelected(tag);
		}
	};

	const onCloseClicked = (e) => {
		e.stopPropagation();
		onTagDeactivated(tag);
	};

	return (
		<Box
			sx={{
				backgroundColor: tag.selected && 'primary.main',
				'&:hover': {
					backgroundColor: 'primary.main',
					opacity: [0.9, 0.8, 0.7],
				},
				cursor: 'pointer',
				borderRadius: '7px',
				padding: '0px 5px 0px 0px',
				display: 'flex',
				flexDirection: 'row',
				alignItems: 'center',
			}}
			onClick={(e) => onTagClicked(e)}
		>
			<IconButton onClick={(e) => onCloseClicked(e)}>
				<CloseIcon />
			</IconButton>
			<Typography variant="button">{tag.title}</Typography>
		</Box>
	);
}

export default ActiveTag;
