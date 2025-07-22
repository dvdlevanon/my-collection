import { useRef } from 'react';

function useVideoController() {
	const videoElement = useRef();

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
		if (videoElement.current) {
			videoElement.current.requestFullscreen();
		}
	};

	const exitFullScreen = () => {
		document.exitFullScreen();
	};

	const getVolume = () => {
		if (videoElement.current) {
			return videoElement.current.volume;
		} else {
			return 0;
		}
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
		getVolume,
		setVolume,

		videoElement,
	};
}

export default useVideoController;
