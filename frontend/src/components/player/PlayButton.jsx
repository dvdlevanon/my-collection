import { useTheme } from '@emotion/react';
import PauseIcon from '@mui/icons-material/Pause';
import PlayIcon from '@mui/icons-material/PlayArrow';
import { IconButton, Tooltip } from '@mui/material';
import { usePlayerStore } from './PlayerStore';

function PlayButton() {
	const theme = useTheme();
	const playerStore = usePlayerStore();

	return (
		<Tooltip title={playerStore.isPlaying ? 'Pause' : 'Play'}>
			<IconButton onClick={playerStore.togglePlay}>
				{playerStore.isPlaying ? (
					<PauseIcon sx={{ fontSize: theme.iconSize(1) }} />
				) : (
					<PlayIcon sx={{ fontSize: theme.iconSize(1) }} />
				)}
			</IconButton>
		</Tooltip>
	);
}

export default PlayButton;
