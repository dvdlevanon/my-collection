import { Stack } from '@mui/material';
import useSize from '@react-hook/size';
import { useQuery } from '@tanstack/react-query';
import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import ReactQueryUtil from '../../utils/react-query-util';
import CropFrame from './CropFrame';
import HighlightControls from './HighlightControls';
import ItemSuggestions from './ItemSuggestions';
import PlayerControls from './PlayerControls';
import { usePlayerStore } from './PlayerStore';
import useVideoController from './VideoController';
import VideoElement from './VideoElement';

function Player({ itemId }) {
	const videoController = useVideoController();
	const playerStore = usePlayerStore();
	const [playerWidth] = useSize(videoController.videoElement);
	const itemQuery = useQuery(ReactQueryUtil.itemQuery(itemId));
	const suggestedQuery = useQuery(ReactQueryUtil.suggestionQuery(itemId));
	const navigate = useNavigate();

	useEffect(() => {
		playerStore.setItemId(itemId);
		playerStore.setNavigate(navigate);
	}, []);

	useEffect(() => {
		let startPosition = itemQuery.data.start_position || 0;
		let initialEndPosition = itemQuery.data.end_position || 0;
		let alreadySplit = itemQuery.data.sub_items;
		playerStore.setUrl(itemQuery.data.url);
		playerStore.setStartTime(startPosition);
		playerStore.setCurrentTime(startPosition);
		playerStore.setEndTime(initialEndPosition);
		playerStore.setDuration(initialEndPosition - startPosition);
		playerStore.setAllowToSplit(!alreadySplit);
	}, [itemQuery.data]);

	useEffect(() => {
		playerStore.setSuggestions(suggestedQuery.data);
	}, [suggestedQuery.data]);

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

	return (
		<Stack
			ref={videoController.videoContainer}
			display="flex"
			sx={{
				position: 'relative',
				cursor: playerStore.isPlaying && !playerStore.controlsVisible ? 'none' : 'auto',
			}}
			tabIndex="0"
			onMouseLeave={() => onMouseLeave()}
		>
			<VideoElement videoController={videoController} />
			<ItemSuggestions width={playerWidth} />
			<HighlightControls />
			<CropFrame videoRef={videoController.videoElement} isPlaying={playerStore.isPlaying} />
			<PlayerControls />
		</Stack>
	);
}

export default Player;
