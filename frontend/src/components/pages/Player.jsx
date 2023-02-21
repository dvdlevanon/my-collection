import FullscreenIcon from '@mui/icons-material/Fullscreen';
import FullscreenExitIcon from '@mui/icons-material/FullscreenExit';
import ImageIcon from '@mui/icons-material/Image';
import LeftIcon from '@mui/icons-material/NavigateBefore';
import RightIcon from '@mui/icons-material/NavigateNext';
import PauseIcon from '@mui/icons-material/Pause';
import PlayIcon from '@mui/icons-material/PlayArrow';
import TimingIcon from '@mui/icons-material/Schedule';
import VolumeDownIcon from '@mui/icons-material/VolumeDown';
import VolumeMuteIcon from '@mui/icons-material/VolumeOff';
import VolumeUpIcon from '@mui/icons-material/VolumeUp';
import { Box, Fade, IconButton, Slider, Stack, Tooltip, Typography } from '@mui/material';
import React, { useEffect, useRef, useState } from 'react';
import Client from '../../network/client';

function Player({ url, setMainCover }) {
	let [showControls, setShowControls] = useState(true);
	let [isPlaying, setIsPlaying] = useState(false);
	let [showVolume, setShowVolume] = useState(false);
	let [showSchedule, setShowSchedule] = useState(false);
	let [volume, setVolume] = useState(0);
	let [currentTime, setCurrentTime] = useState(0);
	let [fullScreen, setFullScreen] = useState(false);
	let [duration, setDuration] = useState(0);
	let videoElement = useRef();
	let playerElement = useRef();

	useEffect(() => {
		volume = parseFloat(localStorage.getItem('volume') || 0.3);
		changeVolume(volume);
	}, []);

	const togglePlay = (e) => {
		if (isPlaying) {
			videoElement.current.pause();
			setIsPlaying(false);
		} else {
			videoElement.current.play();
			setIsPlaying(true);
		}
	};

	const toggleMute = (e) => {
		if (videoElement.current.volume == 0) {
			changeVolume(0.3);
		} else {
			changeVolume(0);
		}

		setVolume(videoElement.current.volume);
	};

	const changeVolume = (volume) => {
		videoElement.current.volume = volume;
		setVolume(videoElement.current.volume);
		localStorage.setItem('volume', volume);
	};

	const changeTime = (newValue) => {
		videoElement.current.currentTime = newValue;
	};

	const formatSeconds = (seconds) => {
		return (
			Math.floor(seconds / 60)
				.toString()
				.padStart(2, '0') +
			':' +
			Math.floor(seconds % 60)
				.toString()
				.padStart(2, '0')
		);
	};

	const enterFullScreen = () => {
		playerElement.current.requestFullscreen();
		setFullScreen(true);
	};

	const exitFullScreen = () => {
		document.exitFullscreen();
		setFullScreen(false);
	};

	const setRelativeTime = (offset) => {
		videoElement.current.currentTime = videoElement.current.currentTime + offset;
	};

	return (
		<Box
			display="flex"
			height="100%"
			width="100%"
			ref={playerElement}
			sx={{
				position: 'relative',
			}}
			onMouseEnter={(e) => {
				setShowControls(true);
			}}
			onMouseLeave={(e) => {
				setShowControls(!isPlaying);
			}}
			tabIndex="0"
			onKeyPress={(e) => {
				console.log('KEY PRESS');
			}}
		>
			<Box
				component="video"
				height="100%"
				width="100%"
				playsInline
				autoPlay={false}
				loop={false}
				ref={videoElement}
				onClick={togglePlay}
				onDoubleClick={fullScreen ? exitFullScreen : enterFullScreen}
				onTimeUpdate={(e) => setCurrentTime(e.target.currentTime)}
				onLoadedMetadata={(e) => setDuration(e.target.duration)}
			>
				<source src={Client.buildFileUrl(url)} />
			</Box>
			<Fade in={showControls}>
				<Stack
					sx={{
						position: 'absolute',
						padding: '10px',
						bottom: -1,
						left: 0,
						right: 0,
						flexDirection: 'column',
						background: 'rgb(255,255,255)',
						background: 'linear-gradient(180deg, rgba(255,255,255,0) 0%, rgba(0,0,0,1) 100%)',
					}}
				>
					<Slider
						min={0}
						max={duration}
						value={currentTime}
						valueLabelDisplay="auto"
						valueLabelFormat={(number) => {
							return formatSeconds(currentTime);
						}}
						onChange={(e, newValue) => changeTime(newValue)}
					/>
					<Stack sx={{ flexDirection: 'row' }}>
						<Tooltip title={isPlaying ? 'Pause' : 'Play'}>
							<IconButton onClick={togglePlay}>{isPlaying ? <PauseIcon /> : <PlayIcon />}</IconButton>
						</Tooltip>
						<Stack
							display="flex"
							flexDirection="row"
							alignItems="center"
							gap="20px"
							onMouseEnter={(e) => setShowVolume(true)}
							onMouseLeave={(e) => setShowVolume(false)}
						>
							<Tooltip title={volume == 0 ? 'Unmute' : 'Mute'}>
								<IconButton onClick={toggleMute}>
									{volume == 0 ? (
										<VolumeMuteIcon />
									) : volume < 0.5 ? (
										<VolumeDownIcon />
									) : (
										<VolumeUpIcon />
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
							<Box>
								<Typography>
									{formatSeconds(currentTime)} / {formatSeconds(duration)}
								</Typography>
							</Box>
						</Stack>
						<Box display="flex" flexGrow={1} justifyContent="flex-end">
							<Stack
								display="flex"
								flexDirection="row"
								alignItems="center"
								gap="20px"
								onMouseEnter={(e) => setShowSchedule(true)}
								onMouseLeave={(e) => setShowSchedule(false)}
							>
								{showSchedule && (
									<Fade in={showSchedule}>
										<Stack flexDirection="row">
											<Tooltip title="-50ms">
												<IconButton size="small" onMouseDown={() => setRelativeTime(-0.005)}>
													<LeftIcon sx={{ fontSize: '20px' }} />
												</IconButton>
											</Tooltip>
											<Tooltip title="-250ms">
												<IconButton size="small" onClick={() => setRelativeTime(-0.025)}>
													<LeftIcon sx={{ fontSize: '25px' }} />
												</IconButton>
											</Tooltip>
											<Tooltip title="-1000ms">
												<IconButton size="small" onClick={() => setRelativeTime(-1.0)}>
													<LeftIcon sx={{ fontSize: '35px' }} />
												</IconButton>
											</Tooltip>
											<Tooltip title="+1000ms">
												<IconButton size="small" onClick={() => setRelativeTime(1.0)}>
													<RightIcon sx={{ fontSize: '35px' }} />
												</IconButton>
											</Tooltip>
											<Tooltip title="+250ms">
												<IconButton size="small" onClick={() => setRelativeTime(0.025)}>
													<RightIcon sx={{ fontSize: '25px' }} />
												</IconButton>
											</Tooltip>
											<Tooltip title="+50ms">
												<IconButton size="small" onClick={() => setRelativeTime(0.005)}>
													<RightIcon sx={{ fontSize: '20px' }} />
												</IconButton>
											</Tooltip>
										</Stack>
									</Fade>
								)}
								<Tooltip title={'Timing controls'}>
									<IconButton>
										<TimingIcon />
									</IconButton>
								</Tooltip>
							</Stack>
							<IconButton onClick={() => setMainCover(videoElement.current.currentTime)}>
								<ImageIcon />
							</IconButton>
							{(!fullScreen && (
								<Tooltip title="Full screen">
									<IconButton onClick={enterFullScreen}>
										<FullscreenIcon />
									</IconButton>
								</Tooltip>
							)) || (
								<Tooltip title="Exit full screen">
									<IconButton onClick={exitFullScreen}>
										<FullscreenExitIcon />
									</IconButton>
								</Tooltip>
							)}
						</Box>
					</Stack>
				</Stack>
			</Fade>
		</Box>
	);
}

export default Player;
