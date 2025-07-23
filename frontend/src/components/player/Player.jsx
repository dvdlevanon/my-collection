import { useTheme } from '@emotion/react';
import { Box, Stack } from '@mui/material';
import { useEffect, useLayoutEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Client from '../../utils/client';
import CropFrame from './CropFrame';
import HighlightControls from './HighlightControls';
import ItemSuggestions from './ItemSuggestions';
import { usePlayerActionStore } from './PlayerActionStore';
import PlayerControls from './PlayerControls';
import { usePlayerStore } from './PlayerStore';
import useVideoController from './VideoController';

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
	const videoController = useVideoController();
	const playerStore = usePlayerStore();
	const playerActionStore = usePlayerActionStore();

	const [playerWidth, setPlayerWidth] = useState(0);
	const [playerHeight, setPlayerHeight] = useState(0);
	const playerElement = useRef();
	const theme = useTheme();
	const navigate = useNavigate();

	useLayoutEffect(() => {
		function updateSize() {
			setPlayerWidth(videoController.videoElement.current.offsetWidth);
			setPlayerHeight(videoController.videoElement.current.offsetHeight);
		}
		window.addEventListener('resize', updateSize);
		const resizeObserver = new ResizeObserver(updateSize);
		resizeObserver.observe(videoController.videoElement.current);

		updateSize();

		return () => {
			window.removeEventListener('resize', updateSize);
			resizeObserver.disconnect();
		};
	}, [videoController.videoElement]);

	useEffect(() => {
		playerStore.setSuggestedItems(suggestedItems);
		playerStore.setNavigate(navigate);
		playerStore.setStartTime(startPosition);
		playerStore.setCurrentTime(startPosition);
		playerStore.setEndTime(initialEndPosition);
		playerStore.setDuration(initialEndPosition - startPosition);
	}, []);

	useEffect(() => {
		playerStore.setVideoController(videoController);
		playerStore.loadFromLocalStorage();
		playerStore.setIsPlaying(true);
	}, [videoController.videoElement]);

	useEffect(() => {
		window.addEventListener('keyup', onKeyDown, false);
		return () => {
			window.removeEventListener('keyup', onKeyDown, false);
		};
	}, [playerStore.isPlaying]);

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
			playerStore.togglePlay();
		} else if (e.key == 'ArrowLeft') {
			playerStore.offsetSeek(e.ctrlKey ? -60 : -10);
		} else if (e.key == 'ArrowRight') {
			playerStore.offsetSeek(e.ctrlKey ? 60 : 10);
		}
	};

	const onMouseLeave = () => {
		if (playerStore.isPlaying) {
			playerStore.hideControls();
		}

		playerStore.setShowVolume(false);
		playerStore.setShowSchedule(false);
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
				ref={videoController.videoElement}
				onClick={playerStore.togglePlay}
				onEnded={playerStore.videoFinished}
				onDoubleClick={playerStore.toggleFullScreen}
				onTimeUpdate={(e) => {
					playerStore.videoTimeUpdate(e.target.currentTime);
				}}
				onLoadedMetadata={(e) => {
					playerStore.videoLoadedMetadata(e.target.duration);
				}}
				onMouseMove={() => playerStore.showControls(true)}
			>
				<source src={Client.buildFileUrl(url)} />
			</Box>
		);
	};

	return (
		<Stack
			display="flex"
			ref={playerElement}
			sx={{
				position: 'relative',
				cursor: playerStore.isPlaying && !playerStore.controlsVisible ? 'none' : 'auto',
			}}
			tabIndex="0"
			onMouseLeave={() => onMouseLeave()}
		>
			{playerActionStore.cropActive() && (
				<CropFrame
					videoRef={videoController.videoElement}
					isPlaying={playerStore.isPlaying}
					width={playerWidth}
					height={playerHeight}
					onMouseMove={() => playerStore.showControls(true)}
				/>
			)}
			{getVideoElement()}
			{playerStore.showSuggestions && suggestedItems && (
				<ItemSuggestions
					suggestedItems={suggestedItems}
					width={playerWidth}
					onBackgroundClick={playerStore.togglePlay}
					onBackgroundDoubleClick={playerStore.toggleFullScreen}
				/>
			)}
			{playerActionStore.highlightActive() && (
				<HighlightControls
					onCancel={playerActionStore.highlightCanceled}
					onDone={(highlightId) => {
						makeHighlight(playerActionStore.startHighlightSecond, playerStore.currentTime, highlightId);
						playerActionStore.highlightCompleted();
					}}
				/>
			)}
			<PlayerControls
				setMainCover={setMainCover}
				splitVideo={splitVideo}
				allowToSplit={allowToSplit}
				cropFrame={cropFrame}
			/>
		</Stack>
	);
}

export default Player;
