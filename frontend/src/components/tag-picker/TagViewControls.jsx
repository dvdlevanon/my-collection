import ZoomInIcon from '@mui/icons-material/ZoomIn';
import ZoomOutIcon from '@mui/icons-material/ZoomOut';
import { IconButton, Stack } from '@mui/material';

function TagViewControls({ tagSize, onZoomChanged }) {
	return (
		<Stack flexDirection="row" gap="10px">
			<IconButton disabled={tagSize <= 100} onClick={() => onZoomChanged(-50)}>
				<ZoomOutIcon />
			</IconButton>
			<IconButton onClick={() => onZoomChanged(50)}>
				<ZoomInIcon />
			</IconButton>
		</Stack>
	);
}

export default TagViewControls;
