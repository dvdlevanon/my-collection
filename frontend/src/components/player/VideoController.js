import { useRef } from 'react';

function useVideoController() {
	const videoElement = useRef();
	const videoContainer = useRef();

	const play = () => {
		if (videoElement.current) {
			videoElement.current.play();
		}
	};

	const pause = () => {
		if (videoElement.current) {
			videoElement.current.pause();
		}
	};

	const seek = (time) => {
		if (videoElement.current) {
			videoElement.current.currentTime = time;
		}
	};

	const currentTime = () => {
		if (videoElement.current) {
			return videoElement.current.currentTime;
		} else {
			return 0;
		}
	};

	const enterFullScreen = () => {
		if (videoContainer.current) {
			videoContainer.current.requestFullscreen();
		}
	};

	const exitFullScreen = () => {
		document.exitFullscreen();
	};

	const setVolume = (volume) => {
		if (videoElement.current) {
			videoElement.current.volume = volume;
		}
	};

	return {
		play,
		pause,
		seek,
		currentTime,
		enterFullScreen,
		exitFullScreen,
		setVolume,

		videoElement,
		videoContainer,
	};
}

export default useVideoController;
