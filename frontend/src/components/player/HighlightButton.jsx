import { useTheme } from '@emotion/react';
import HighlightIcon from '@mui/icons-material/Stars';
import { IconButton } from '@mui/material';
import { usePlayerActionStore } from './PlayerActionStore';
import { usePlayerStore } from './PlayerStore';

function HighlightButton() {
	const theme = useTheme();
	const playerStore = usePlayerStore();
	const playerActionStore = usePlayerActionStore();

	return (
		<IconButton onClick={() => playerActionStore.startHighlightCreation(playerStore.currentTime)}>
			<HighlightIcon sx={{ fontSize: theme.iconSize(1) }} />
		</IconButton>
	);
}

export default HighlightButton;
