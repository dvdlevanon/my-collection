import { useTheme } from '@emotion/react';
import VolumeDownIcon from '@mui/icons-material/VolumeDown';
import VolumeMuteIcon from '@mui/icons-material/VolumeOff';
import VolumeUpIcon from '@mui/icons-material/VolumeUp';
import { Fade, IconButton, Slider, Stack, Tooltip } from '@mui/material';
import { usePlayerStore } from './PlayerStore';

function VolumeControls() {
	const theme = useTheme();
	const playerStore = usePlayerStore();

	const toggleMute = (e) => {
		if (playerStore.volume == 0) {
			changeVolume(0.3);
		} else {
			changeVolume(0);
		}
	};

	const changeVolume = (volume) => {
		playerStore.setVolume(volume);
	};

	return (
		<Stack
			display="flex"
			flexDirection="row"
			alignItems="center"
			gap={theme.spacing(2)}
			onMouseEnter={(e) => playerStore.setShowVolume(true)}
		>
			<Tooltip title={playerStore.volume == 0 ? 'Unmute' : 'Mute'}>
				<IconButton onClick={toggleMute}>
					{playerStore.volume == 0 ? (
						<VolumeMuteIcon sx={{ fontSize: theme.iconSize(1) }} />
					) : playerStore.volume < 0.5 ? (
						<VolumeDownIcon sx={{ fontSize: theme.iconSize(1) }} />
					) : (
						<VolumeUpIcon sx={{ fontSize: theme.iconSize(1) }} />
					)}
				</IconButton>
			</Tooltip>
			{playerStore.showVolume && (
				<Fade in={playerStore.showVolume}>
					<Slider
						min={0}
						max={100}
						value={playerStore.volume * 100}
						onChange={(e, newValue) => changeVolume(newValue / 100)}
						sx={{
							width: '100px',
						}}
					/>
				</Fade>
			)}
		</Stack>
	);
}

export default VolumeControls;
