import { useTheme } from '@emotion/react';
import CancelIcon from '@mui/icons-material/Cancel';
import SplitIcon from '@mui/icons-material/ContentCut';
import CropIcon from '@mui/icons-material/Crop';
import DoneIcon from '@mui/icons-material/Done';
import FullscreenIcon from '@mui/icons-material/Fullscreen';
import FullscreenExitIcon from '@mui/icons-material/FullscreenExit';
import ImageIcon from '@mui/icons-material/Image';
import PauseIcon from '@mui/icons-material/Pause';
import PlayIcon from '@mui/icons-material/PlayArrow';
import PlayNextIcon from '@mui/icons-material/QueuePlayNext';
import HighlightIcon from '@mui/icons-material/Stars';
import { Box, Fade, IconButton, Stack, Tooltip, Typography } from '@mui/material';
import { useEffect, useLayoutEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Client from '../../utils/client';
import TimeUtil from '../../utils/time-utils';
import CropFrame from './CropFrame';
import HighlightControls from './HighlightControls';
import ItemSuggestions from './ItemSuggestions';
import PlayerSlider from './PlayerSlider';
import TimingControls from './TimingControls';
import VolumeControls from './VolumeControls';

function Player({
	url,
	setMainCover,
	startPosition,
	initialEndPosition,
	splitVideo,
	makeHighlight,
	cropFrame,
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
	const [autoPlayNext, setAutoPlayNext] = useState(false);
	const [duration, setDuration] = useState(initialEndPosition - startPosition);
	const [hideControlsTimerId, setHideControlsTimerId] = useState(0);
	const [playerWidth, setPlayerWidth] = useState(0);
	const [playerHeight, setPlayerHeight] = useState(0);
	const [startHighlightSecond, setStartHighlightSecond] = useState(-1);
	const [cropMode, setCropMode] = useState(false);
	const [frameCrop, setFrameCrop] = useState(null);
	const videoElement = useRef();
	const playerElement = useRef();
	const theme = useTheme();
	const navigate = useNavigate();

	useLayoutEffect(() => {
		function updateSize() {
			setPlayerWidth(videoElement.current.offsetWidth);
			setPlayerHeight(videoElement.current.offsetHeight);
		}
		window.addEventListener('resize', updateSize);
		const resizeObserver = new ResizeObserver(updateSize);
		resizeObserver.observe(videoElement.current);

		updateSize();

		return () => {
			window.removeEventListener('resize', updateSize);
			resizeObserver.disconnect();
		};
	}, [videoElement]);

	useEffect(() => {
		let autoPlayNext = localStorage.getItem('auto-play-next');

		if (autoPlayNext) {
			setAutoPlayNext(autoPlayNext == 'true');
		}
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

	const enterFullScreen = () => {
		playerElement.current.requestFullscreen();
		setFullScreen(true);
	};

	const exitFullScreen = () => {
		document.exitFullscreen();
		setFullScreen(false);
	};

	const videoFinished = () => {
		if (autoPlayNext && suggestedItems) {
			let nextItemIndex = Math.floor(Math.random() * suggestedItems.length);
			navigate('/spa/item/' + suggestedItems[nextItemIndex].id);
		} else {
			setShowSuggestions(true);
		}

		setIsPlaying(false);
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

	const getVideoElement = () => {
		return (
			<Box
				borderRadius={theme.spacing(2)}
				sx={{
					boxShadow: '3',
				}}
				component="video"
				crossOrigin="anonymous"
				height="100%"
				width="100%"
				playsInline
				autoPlay={true}
				loop={false}
				ref={videoElement}
				onClick={togglePlay}
				onEnded={() => {
					videoFinished();
				}}
				onDoubleClick={fullScreen ? exitFullScreen : enterFullScreen}
				onTimeUpdate={(e) => {
					setCurrentTime(e.target.currentTime);
					if (endPosition > 0 && e.target.currentTime >= endPosition) {
						e.target.currentTime = startPosition;
						videoElement.current.pause();
						videoFinished();
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
		);
	};

	const cropClicked = () => {
		videoElement.current.pause();
		setIsPlaying(false);
		setCropMode(true);
	};

	const cancelCropClicked = () => {
		setCropMode(false);
	};

	const finishCropClicked = () => {
		setCropMode(false);
		cropFrame(currentTime, frameCrop);
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
			{cropMode && (
				<CropFrame
					videoRef={videoElement}
					isPlaying={isPlaying}
					width={playerWidth}
					height={playerHeight}
					onMouseMove={() => onMouseMove()}
					setCrop={setFrameCrop}
				/>
			)}
			{getVideoElement()}
			{showSuggestions && suggestedItems && (
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
					zIndex={2}
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
					<PlayerSlider
						min={startPosition}
						max={endPosition}
						value={currentTime}
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
						<Tooltip title={'Toggle Auto Play Next'}>
							<IconButton
								onClick={() => {
									let newValue = !autoPlayNext;
									setAutoPlayNext(!autoPlayNext);
									localStorage.setItem('auto-play-next', newValue);
								}}
							>
								{
									<PlayNextIcon
										color={autoPlayNext ? 'secondary' : 'auto'}
										sx={{ fontSize: theme.iconSize(1) }}
									/>
								}
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
								{TimeUtil.formatSeconds(currentTime - startPosition)} /{' '}
								{TimeUtil.formatSeconds(duration)}
							</Typography>
						</Box>
						<Box display="flex" flexGrow={1} justifyContent="flex-end">
							<TimingControls
								setRelativeTime={setRelativeTime}
								showSchedule={showSchedule}
								setShowSchedule={setShowSchedule}
							/>
							{!cropMode && (
								<IconButton onClick={() => setMainCover(videoElement.current.currentTime)}>
									<ImageIcon sx={{ fontSize: theme.iconSize(1) }} />
								</IconButton>
							)}
							{cropMode ? (
								<>
									<IconButton onClick={finishCropClicked}>
										<DoneIcon sx={{ fontSize: theme.iconSize(1) }} />
									</IconButton>
									<IconButton onClick={cancelCropClicked}>
										<CancelIcon sx={{ fontSize: theme.iconSize(1) }} />
									</IconButton>
								</>
							) : (
								<IconButton onClick={cropClicked}>
									<CropIcon sx={{ fontSize: theme.iconSize(1) }} />
								</IconButton>
							)}
							{!cropMode && (
								<IconButton
									disabled={!allowToSplit()}
									onClick={() => splitVideo(videoElement.current.currentTime)}
								>
									<SplitIcon sx={{ fontSize: theme.iconSize(1) }} />
								</IconButton>
							)}
							{!cropMode && (
								<IconButton onClick={() => setStartHighlightSecond(videoElement.current.currentTime)}>
									<HighlightIcon sx={{ fontSize: theme.iconSize(1) }} />
								</IconButton>
							)}
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
