import { useTheme } from '@emotion/react';
import VolumeDownIcon from '@mui/icons-material/VolumeDown';
import VolumeMuteIcon from '@mui/icons-material/VolumeOff';
import VolumeUpIcon from '@mui/icons-material/VolumeUp';
import { Fade, IconButton, Slider, Stack, Tooltip } from '@mui/material';
import { useEffect, useState } from 'react';
import { usePlayerStore } from './PlayerStore';

function VolumeControls() {
	const setShowVolume = usePlayerStore((state) => state.setShowVolume);
	const showVolume = usePlayerStore((state) => state.showVolume);
	const [volume, setVolume] = useState(0);
	const theme = useTheme();
	const playerStore = usePlayerStore();

	useEffect(() => {
		let volume = parseFloat(localStorage.getItem('volume') || 0.3);
		changeVolume(volume);
	}, []);

	const toggleMute = (e) => {
		if (playerStore.getVolume() == 0) {
			changeVolume(0.3);
		} else {
			changeVolume(0);
		}

		setVolume(playerStore.getVolume());
	};

	const changeVolume = (volume) => {
		playerStore.setVolume(volume);
		setVolume(playerStore.getVolume());
		localStorage.setItem('volume', volume);
	};

	return (
		<Stack
			display="flex"
			flexDirection="row"
			alignItems="center"
			gap={theme.spacing(2)}
			onMouseEnter={(e) => setShowVolume(true)}
		>
			<Tooltip title={volume == 0 ? 'Unmute' : 'Mute'}>
				<IconButton onClick={toggleMute}>
					{volume == 0 ? (
						<VolumeMuteIcon sx={{ fontSize: theme.iconSize(1) }} />
					) : volume < 0.5 ? (
						<VolumeDownIcon sx={{ fontSize: theme.iconSize(1) }} />
					) : (
						<VolumeUpIcon sx={{ fontSize: theme.iconSize(1) }} />
					)}
				</IconButton>
			</Tooltip>
			{showVolume && (
				<Fade in={showVolume}>
					<Slider
						min={0}
						max={100}
						value={volume * 100}
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
