import { useTheme } from '@emotion/react';
import SplitIcon from '@mui/icons-material/ContentCut';
import FullscreenIcon from '@mui/icons-material/Fullscreen';
import FullscreenExitIcon from '@mui/icons-material/FullscreenExit';
import ImageIcon from '@mui/icons-material/Image';
import PauseIcon from '@mui/icons-material/Pause';
import PlayIcon from '@mui/icons-material/PlayArrow';
import PlayNextIcon from '@mui/icons-material/QueuePlayNext';
import HighlightIcon from '@mui/icons-material/Stars';
import { Box, Fade, IconButton, Slider, Stack, Tooltip, Typography } from '@mui/material';
import React, { useEffect, useLayoutEffect, useRef, useState } from 'react';
import Client from '../../utils/client';
import HighlightControls from './HighlightControls';
import ItemSuggestions from './ItemSuggestions';
import TimingControls from './TimingControls';
import VolumeControls from './VolumeControls';

function Player({
	url,
	setMainCover,
	startPosition,
	initialEndPosition,
	splitVideo,
	makeHighlight,
	allowToSplit,
	suggestedItems,
}) {
	const [showControls, setShowControls] = useState(true);
	const [showVolume, setShowVolume] = useState(false);
	const [showSchedule, setShowSchedule] = useState(false);
	const [isPlaying, setIsPlaying] = useState(true);
	const [currentTime, setCurrentTime] = useState(startPosition);
	const [endPosition, setEndPosition] = useState(initialEndPosition);
	const [fullScreen, setFullScreen] = useState(false);
	const [showSuggestions, setShowSuggestions] = useState(false);
	const [duration, setDuration] = useState(initialEndPosition - startPosition);
	const [hideControlsTimerId, setHideControlsTimerId] = useState(0);
	const [playerWidth, setPlayerWidth] = useState(0);
	const [startHighlightSecond, setStartHighlightSecond] = useState(-1);
	const videoElement = useRef();
	const playerElement = useRef();
	const theme = useTheme();

	useLayoutEffect(() => {
		function updateSize() {
			setPlayerWidth(videoElement.current.offsetWidth);
		}
		window.addEventListener('resize', updateSize);
		updateSize();
		return () => window.removeEventListener('resize', updateSize);
	}, []);

	useEffect(() => {
		window.addEventListener('keyup', onKeyDown, false);
		return () => {
			window.removeEventListener('keyup', onKeyDown, false);
		};
	}, [isPlaying]);

	const isInputFocused = () => {
		var activeElement = document.activeElement;
		var inputs = ['input', 'select', 'button', 'textarea'];
		return activeElement && inputs.indexOf(activeElement.tagName.toLowerCase()) !== -1;
	};

	const onKeyDown = (e) => {
		if (isInputFocused()) {
			return;
		}

		if (e.key == ' ') {
			togglePlay();
		} else if (e.key == 'ArrowLeft') {
			setRelativeTime(e.ctrlKey ? -60 : -10);
		} else if (e.key == 'ArrowRight') {
			setRelativeTime(e.ctrlKey ? 60 : 10);
		}
	};

	const togglePlay = (e) => {
		if (isPlaying) {
			videoElement.current.pause();
			setIsPlaying(false);
		} else {
			videoElement.current.play();
			setIsPlaying(true);
			setShowSuggestions(false);
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
		let newOffset = videoElement.current.currentTime + offset;
		if (newOffset > endPosition) {
			videoElement.current.currentTime = endPosition;
		} else if (newOffset < startPosition) {
			videoElement.current.currentTime = startPosition;
		} else {
			videoElement.current.currentTime = newOffset;
		}
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
		<Stack
			display="flex"
			ref={playerElement}
			sx={{
				position: 'relative',
				cursor: isPlaying && !showControls ? 'none' : 'auto',
			}}
			tabIndex="0"
			onMouseLeave={() => onMouseLeave()}
		>
			<Box
				borderRadius={theme.spacing(2)}
				sx={{
					boxShadow: '3',
				}}
				component="video"
				height="100%"
				width="100%"
				playsInline
				autoPlay={true}
				loop={false}
				ref={videoElement}
				onClick={togglePlay}
				onEnded={() => {
					setShowSuggestions(true);
					setIsPlaying(false);
				}}
				onDoubleClick={fullScreen ? exitFullScreen : enterFullScreen}
				onTimeUpdate={(e) => {
					setCurrentTime(e.target.currentTime);
					if (endPosition > 0 && e.target.currentTime >= endPosition) {
						e.target.currentTime = startPosition;
						videoElement.current.pause();
						setShowSuggestions(true);
						setIsPlaying(false);
					}
				}}
				onLoadedMetadata={(e) => {
					e.target.currentTime = startPosition;
					if (endPosition == 0) {
						setEndPosition(e.target.duration);
						setDuration(e.target.duration);
					}
				}}
				onMouseMove={() => onMouseMove()}
			>
				<source src={Client.buildFileUrl(url)} />
			</Box>
			{showSuggestions && (
				<ItemSuggestions
					suggestedItems={suggestedItems}
					width={playerWidth}
					onBackgroundClick={togglePlay}
					onBackgroundDoubleClick={fullScreen ? exitFullScreen : enterFullScreen}
				/>
			)}
			{startHighlightSecond !== -1 && (
				<HighlightControls
					onCancel={() => setStartHighlightSecond(-1)}
					onDone={(highlightId) => {
						makeHighlight(startHighlightSecond, videoElement.current.currentTime, highlightId);
						setStartHighlightSecond(-1);
					}}
				/>
			)}
			<Fade in={showControls}>
				<Stack
					onMouseEnter={() => onMouseEnter()}
					sx={{
						position: 'absolute',
						padding: theme.spacing(1),
						bottom: -1,
						left: 0,
						right: 0,
						flexDirection: 'column',
						background:
							'linear-gradient(180deg, ' +
							theme.palette.gradient.color1 +
							'00 0%, ' +
							theme.palette.gradient.color2 +
							'FF 100%)',
					}}
				>
					<Slider
						min={startPosition}
						max={endPosition}
						value={currentTime}
						valueLabelDisplay="auto"
						valueLabelFormat={(number) => {
							return formatSeconds(currentTime - startPosition);
						}}
						onChange={(e, newValue) => changeTime(newValue)}
					/>
					<Stack flexDirection="row" alignItems="center" gap={theme.spacing(2)}>
						<Tooltip title={isPlaying ? 'Pause' : 'Play'}>
							<IconButton onClick={togglePlay}>
								{isPlaying ? (
									<PauseIcon sx={{ fontSize: theme.iconSize(1) }} />
								) : (
									<PlayIcon sx={{ fontSize: theme.iconSize(1) }} />
								)}
							</IconButton>
						</Tooltip>
						<Tooltip title={showSuggestions ? 'Hide Suggestions' : 'Show Suggestions'}>
							<IconButton onClick={() => setShowSuggestions(!showSuggestions)}>
								{<PlayNextIcon sx={{ fontSize: theme.iconSize(1) }} />}
							</IconButton>
						</Tooltip>
						<VolumeControls
							showVolume={showVolume}
							setShowVolume={setShowVolume}
							getVideoVolume={() => videoElement.current.volume}
							setVideoVolume={(volume) => (videoElement.current.volume = volume)}
						/>
						<Box>
							<Typography>
								{formatSeconds(currentTime - startPosition)} / {formatSeconds(duration)}
							</Typography>
						</Box>
						<Box display="flex" flexGrow={1} justifyContent="flex-end">
							<TimingControls
								setRelativeTime={setRelativeTime}
								showSchedule={showSchedule}
								setShowSchedule={setShowSchedule}
							/>
							<IconButton onClick={() => setMainCover(videoElement.current.currentTime)}>
								<ImageIcon sx={{ fontSize: theme.iconSize(1) }} />
							</IconButton>
							<IconButton
								disabled={!allowToSplit()}
								onClick={() => splitVideo(videoElement.current.currentTime)}
							>
								<SplitIcon sx={{ fontSize: theme.iconSize(1) }} />
							</IconButton>
							<IconButton onClick={() => setStartHighlightSecond(videoElement.current.currentTime)}>
								<HighlightIcon sx={{ fontSize: theme.iconSize(1) }} />
							</IconButton>
							{(!fullScreen && (
								<Tooltip title="Full screen">
									<IconButton onClick={enterFullScreen}>
										<FullscreenIcon sx={{ fontSize: theme.iconSize(1) }} />
									</IconButton>
								</Tooltip>
							)) || (
								<Tooltip title="Exit full screen">
									<IconButton onClick={exitFullScreen}>
										<FullscreenExitIcon sx={{ fontSize: theme.iconSize(1) }} />
									</IconButton>
								</Tooltip>
							)}
						</Box>
					</Stack>
				</Stack>
			</Fade>
		</Stack>
	);
}

export default Player;
