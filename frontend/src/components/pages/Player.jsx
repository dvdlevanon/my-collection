import FullscreenIcon from '@mui/icons-material/Fullscreen';
import FullscreenExitIcon from '@mui/icons-material/FullscreenExit';
import PauseIcon from '@mui/icons-material/Pause';
import PlayIcon from '@mui/icons-material/PlayArrow';
import VolumeDownIcon from '@mui/icons-material/VolumeDown';
import VolumeMuteIcon from '@mui/icons-material/VolumeOff';
import VolumeUpIcon from '@mui/icons-material/VolumeUp';
import { Box, Fade, IconButton, Slider, Stack, Typography } from '@mui/material';
import React, { useEffect, useLayoutEffect, useRef, useState } from 'react';
import Client from '../../network/client';

function useWindowSize() {
	const [size, setSize] = useState([0, 0]);
	useLayoutEffect(() => {
		function updateSize() {
			setSize([window.innerWidth, window.innerHeight]);
		}
		window.addEventListener('resize', updateSize);
		updateSize();
		return () => window.removeEventListener('resize', updateSize);
	}, []);
	return size;
}

function Player({ url }) {
	let [showControls, setShowControls] = useState(true);
	let [isPlaying, setIsPlaying] = useState(false);
	let [showVolume, setShowVolume] = useState(false);
	let [volume, setVolume] = useState(0);
	let [currentTime, setCurrentTime] = useState(0);
	let [fullScreen, setFullScreen] = useState(false);
	let [duration, setDuration] = useState(0);
	let [videoHeight, setVideoHeight] = useState(0);
	let [videoWidth, setVideoWidth] = useState(0);
	const [windowWidth, windowHeight] = useWindowSize();
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

	const calcControlsOffset = () => {
		if (videoHeight >= videoWidth) {
			return videoWidth;
		}

		let componentWidth = videoElement.current.offsetWidth;
		let componentHeight = videoElement.current.offsetHeight;

		let ratio = videoWidth / videoHeight;
		let actualWidth = componentHeight * ratio;
		let gap = (componentWidth - actualWidth) / 2;
		return gap;
	};

	return (
		<Box
			height={(windowHeight / 5) * 3}
			minHeight={600}
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
				onLoadedMetadata={(e) => {
					setDuration(e.target.duration);
					setVideoWidth(e.target.videoWidth);
					setVideoHeight(e.target.videoHeight);
				}}
			>
				<source src={Client.buildFileUrl(url)} />
			</Box>
			<Fade in={showControls}>
				<Stack
					sx={{
						position: 'absolute',
						bottom: -1,
						left: calcControlsOffset(),
						right: calcControlsOffset(),
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
						<IconButton onClick={togglePlay}>{isPlaying ? <PauseIcon /> : <PlayIcon />}</IconButton>
						<Stack
							display="flex"
							flexDirection="row"
							alignItems="center"
							gap="20px"
							onMouseEnter={(e) => setShowVolume(true)}
							onMouseLeave={(e) => setShowVolume(false)}
						>
							<IconButton onClick={toggleMute}>
								{volume == 0 ? (
									<VolumeMuteIcon />
								) : volume < 0.5 ? (
									<VolumeDownIcon />
								) : (
									<VolumeUpIcon />
								)}
							</IconButton>
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
							{(!fullScreen && (
								<IconButton onClick={enterFullScreen}>
									<FullscreenIcon />
								</IconButton>
							)) || (
								<IconButton onClick={exitFullScreen}>
									<FullscreenExitIcon />
								</IconButton>
							)}
						</Box>
					</Stack>
				</Stack>
			</Fade>
		</Box>
	);
}

export default Player;
