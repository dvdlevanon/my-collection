import { useTheme } from '@emotion/react';
import ZoomInIcon from '@mui/icons-material/ZoomIn';
import ZoomOutIcon from '@mui/icons-material/ZoomOut';
import { IconButton, Stack } from '@mui/material';

function TagViewControls({ tagSize, onZoomChanged }) {
	const theme = useTheme();

	return (
		<Stack flexDirection="row" gap={theme.spacing(1)}>
			<IconButton disabled={tagSize <= 100} onClick={() => onZoomChanged(-50)}>
				<ZoomOutIcon sx={{ fontSize: theme.iconSize(1) }} />
			</IconButton>
			<IconButton onClick={() => onZoomChanged(50)}>
				<ZoomInIcon sx={{ fontSize: theme.iconSize(1) }} />
			</IconButton>
		</Stack>
	);
}

export default TagViewControls;
