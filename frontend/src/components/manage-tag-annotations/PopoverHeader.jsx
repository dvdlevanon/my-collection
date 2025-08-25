import CloseIcon from '@mui/icons-material/Close';
import { Box, IconButton, Typography, useTheme } from '@mui/material';
import { useTagAnnotationsStore } from './tagAnnotationsStore';

function PopoverHeader({ onClose }) {
	const theme = useTheme();
	const tagAnnotationsStore = useTagAnnotationsStore();

	return (
		<Box
			onClick={(e) => {
				e.preventDefault();
				e.stopPropagation();
			}}
			sx={{
				display: 'flex',
				gap: theme.spacing(1),
				alignItems: 'center',
			}}
		>
			<IconButton
				onClick={(e) => {
					e.preventDefault();
					e.stopPropagation();
					onClose(e);
				}}
			>
				<CloseIcon sx={{ fontSize: theme.iconSize(1) }} />
			</IconButton>
			<Typography
				variant="body1"
				noWrap
				onClick={(e) => {
					e.preventDefault();
					e.stopPropagation();
				}}
			>
				{tagAnnotationsStore.tag.title} Annotations
			</Typography>
		</Box>
	);
}

export default PopoverHeader;
