import { create } from 'zustand';

export const usePlayerStore = create((set, get) => ({
	itemId: false,
	url: '',
	controlsVisible: false,
	showVolume: false,
	showSchedule: false,
	showSuggestions: false,
	suggestions: [],
	fullScreen: false,
	autoPlayNext: false,
	allowToSplit: false,
	videoController: null,
	navigate: null,
	suggestedItems: null,
	isPlaying: false,
	startTime: 0,
	endTime: 0,
	currentTime: 0,
	duration: 0,
	hideControlsTimerId: 0,
	volume: 0,

	setSuggestions: (suggestions) => set({ suggestions }),
	setUrl: (url) => set({ url }),
	setItemId: (itemId) => set({ itemId }),
	setShowVolume: (showVolume) => set({ showVolume }),
	setShowSchedule: (showSchedule) => set({ showSchedule }),
	setShowSuggestions: (showSuggestions) => set({ showSuggestions }),
	setVideoController: (videoController) => set({ videoController }),
	setNavigate: (navigate) => set({ navigate }),
	setAllowToSplit: (allowToSplit) => set({ allowToSplit }),
	setSuggestedItems: (suggestedItems) => set({ suggestedItems }),
	setIsPlaying: (isPlaying) => set({ isPlaying }),
	setStartTime: (startTime) => set({ startTime }),
	setEndTime: (endTime) => set({ endTime }),
	setCurrentTime: (currentTime) => set({ currentTime }),
	setDuration: (duration) => set({ duration }),

	loadFromLocalStorage: () => {
		const { setAutoPlayNext, setVolume } = get();

		setAutoPlayNext(localStorage.getItem('auto-play-next') == 'true');
		setVolume(parseFloat(localStorage.getItem('volume') || 0.3));
	},

	setAutoPlayNext: (autoPlayNext) => {
		set({ autoPlayNext });
		localStorage.setItem('auto-play-next', autoPlayNext);
	},

	enterFullScreen: () => {
		const { videoController, setFullScreen } = get();

		videoController.enterFullScreen();
		set({ fullscreen: true });
	},

	exitFullScreen: () => {
		const { videoController } = get();

		videoController.exitFullScreen();
		set({ fullscreen: false });
	},

	toggleFullScreen: () => {
		const { fullScreen, exitFullScreen, enterFullScreen } = get();

		if (fullScreen) {
			exitFullScreen();
		} else {
			enterFullScreen();
		}
	},

	play: () => {
		const { videoController, setIsPlaying, setShowSuggestions } = get();

		videoController.play();
		setIsPlaying(true);
		setShowSuggestions(false);
	},

	pause: () => {
		const { videoController, setIsPlaying } = get();

		videoController.pause();
		setIsPlaying(false);
	},

	togglePlay: () => {
		const { play, pause, isPlaying } = get();

		if (isPlaying) {
			pause();
		} else {
			play();
		}
	},

	videoLoadedMetadata: (duration) => {
		const { seek, startTime, endTime, setEndTime, setDuration } = get();

		seek(startTime);
		if (endTime == 0) {
			setEndTime(duration);
			setDuration(duration);
		}
	},

	videoTimeUpdate: (time) => {
		const { videoController, videoFinished, setCurrentTime, startTime, endTime } = get();

		setCurrentTime(time);
		if (endTime > 0 && time >= endTime) {
			videoController.seek(startTime);
			videoController.pause();
			videoFinished();
		}
	},

	videoFinished: () => {
		const { autoPlayNext, suggestedItems, navigate, setShowSuggestions, setIsPlaying } = get();

		if (autoPlayNext && suggestedItems) {
			let nextItemIndex = Math.floor(Math.random() * suggestedItems.length);
			navigate('/spa/item/' + suggestedItems[nextItemIndex].id);
		} else {
			setShowSuggestions(true);
		}

		setIsPlaying(false);
	},

	seek: (time) => {
		const { videoController } = get();

		videoController.seek(time);
	},

	offsetSeek: (offset) => {
		const { videoController, startTime, endTime, seek } = get();

		let newTime = videoController.currentTime() + offset;
		if (newTime > endTime) {
			seek(endTime);
		} else if (newTime < startTime) {
			seek(startTime);
		} else {
			seek(newTime);
		}
	},

	setVolume: (volume) => {
		const { videoController } = get();
		set({ volume });
		localStorage.setItem('volume', volume);

		if (videoController) {
			videoController.setVolume(volume);
		}
	},

	showControls: (autoHide) => {
		const { hideControlsTimerId, isPlaying } = get();

		if (hideControlsTimerId > 0) {
			clearTimeout(hideControlsTimerId);
		}

		set({ controlsVisible: true, hideControlsTimerId: 0 });

		if (isPlaying && autoHide) {
			const timerId = setTimeout(() => {
				set({ controlsVisible: false, hideControlsTimerId: 0 });
			}, 2000);

			set({ hideControlsTimerId: timerId });
		}
	},

	hideControls: () => {
		const { hideControlsTimerId } = get();

		if (hideControlsTimerId > 0) {
			clearTimeout(hideControlsTimerId);
		}

		set({ controlsVisible: false, hideControlsTimerId: 0 });
	},
}));
