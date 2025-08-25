import { useTheme } from '@emotion/react';
import FullscreenIcon from '@mui/icons-material/Fullscreen';
import FullscreenExitIcon from '@mui/icons-material/FullscreenExit';
import { IconButton, Tooltip } from '@mui/material';
import { usePlayerStore } from './PlayerStore';

function FullscreenButton() {
	const theme = useTheme();
	const playerStore = usePlayerStore();

	return (
		(!playerStore.fullScreen && (
			<Tooltip title="Full screen">
				<IconButton onClick={playerStore.enterFullScreen}>
					<FullscreenIcon sx={{ fontSize: theme.iconSize(1) }} />
				</IconButton>
			</Tooltip>
		)) || (
			<Tooltip title="Exit full screen">
				<IconButton onClick={playerStore.exitFullScreen}>
					<FullscreenExitIcon sx={{ fontSize: theme.iconSize(1) }} />
				</IconButton>
			</Tooltip>
		)
	);
}

export default FullscreenButton;
