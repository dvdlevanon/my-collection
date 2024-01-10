import { useTheme } from '@emotion/react';
import VolumeDownIcon from '@mui/icons-material/VolumeDown';
import VolumeMuteIcon from '@mui/icons-material/VolumeOff';
import VolumeUpIcon from '@mui/icons-material/VolumeUp';
import { Fade, IconButton, Slider, Stack, Tooltip } from '@mui/material';
import React, { useEffect, useState } from 'react';

function VolumeControls({ showVolume, setShowVolume, getVideoVolume, setVideoVolume }) {
	const [volume, setVolume] = useState(0);
	const theme = useTheme();

	useEffect(() => {
		let volume = parseFloat(localStorage.getItem('volume') || 0.3);
		changeVolume(volume);
	}, []);

	const toggleMute = (e) => {
		if (getVideoVolume() == 0) {
			changeVolume(0.3);
		} else {
			changeVolume(0);
		}

		setVolume(getVideoVolume());
	};

	const changeVolume = (volume) => {
		setVideoVolume(volume);
		setVolume(getVideoVolume());
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
