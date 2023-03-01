import FullscreenIcon from '@mui/icons-material/Fullscreen';
import FullscreenExitIcon from '@mui/icons-material/FullscreenExit';
import ImageIcon from '@mui/icons-material/Image';
import PauseIcon from '@mui/icons-material/Pause';
import PlayIcon from '@mui/icons-material/PlayArrow';
import { Box, Fade, IconButton, Slider, Stack, Tooltip, Typography } from '@mui/material';
import React, { useRef, useState } from 'react';
import Client from '../../utils/client';
import TimingControls from './TimingControls';
import VolumeControls from './VolumeControls';

function Player({ url, setMainCover }) {
	let [showControls, setShowControls] = useState(true);
	let [showVolume, setShowVolume] = useState(false);
	let [showSchedule, setShowSchedule] = useState(false);
	let [isPlaying, setIsPlaying] = useState(false);
	let [currentTime, setCurrentTime] = useState(0);
	let [fullScreen, setFullScreen] = useState(false);
	let [duration, setDuration] = useState(0);
	let [hideControlsTimerId, setHideControlsTimerId] = useState(0);
	let videoElement = useRef();
	let playerElement = useRef();

	const togglePlay = (e) => {
		if (isPlaying) {
			videoElement.current.pause();
			setIsPlaying(false);
		} else {
			videoElement.current.play();
			setIsPlaying(true);
		}
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

	const onMouseEnter = () => {
		setShowControls(true);

		if (hideControlsTimerId > 0) {
			clearTimeout(hideControlsTimerId);
		}
	};

	const onMouseLeave = () => {
		setShowControls(!isPlaying);
		setShowVolume(false);
		setShowSchedule(false);
	};

	const onMouseMove = () => {
		onMouseEnter();

		if (isPlaying) {
			setHideControlsTimerId(
				setTimeout(() => {
					setShowControls(false);
				}, 2000)
			);
		}
	};

	return (
		<Box
			display="flex"
			height="100%"
			width="100%"
			ref={playerElement}
			sx={{
				position: 'relative',
				cursor: isPlaying && !showControls ? 'none' : 'auto',
			}}
			tabIndex="0"
			onKeyPress={(e) => {
				console.log('KEY PRESS');
			}}
			onMouseLeave={() => onMouseLeave()}
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
				onMouseMove={() => onMouseMove()}
			>
				<source src={Client.buildFileUrl(url)} />
			</Box>
			<Fade in={showControls}>
				<Stack
					onMouseEnter={() => onMouseEnter()}
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
					<Stack
						sx={{
							flexDirection: 'row',
							alignItems: 'center',
							gap: '20px',
						}}
					>
						<Tooltip title={isPlaying ? 'Pause' : 'Play'}>
							<IconButton onClick={togglePlay}>{isPlaying ? <PauseIcon /> : <PlayIcon />}</IconButton>
						</Tooltip>
						<VolumeControls
							showVolume={showVolume}
							setShowVolume={setShowVolume}
							getVideoVolume={() => videoElement.current.volume}
							setVideoVolume={(volume) => (videoElement.current.volume = volume)}
						/>
						<Box>
							<Typography>
								{formatSeconds(currentTime)} / {formatSeconds(duration)}
							</Typography>
						</Box>
						<Box display="flex" flexGrow={1} justifyContent="flex-end">
							<TimingControls
								setRelativeTime={setRelativeTime}
								showSchedule={showSchedule}
								setShowSchedule={setShowSchedule}
							/>
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
