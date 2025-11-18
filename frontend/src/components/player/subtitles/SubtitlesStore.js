import { create } from 'zustand';

export const useSubtitleStore = create((set, get) => ({
	controlsShown: false,
	fontSize: 4.44,
	fontShadowWidth: 1,
	fontColor: '#ffffff',
	fontShadowColor: '#000000',
	subtitleOffsetMillis: 0,
	selectedSubtitleUrl: '',

	loadFromLocalStorage: () => {
		const { setFontColor, setFontShadowColor, setFontSize, setFontShadowWidth } = get();

		setFontColor(localStorage.getItem('subtitles-font-color') || '#ffffff');
		setFontShadowColor(localStorage.getItem('subtitles-font-shadow-color') || '#000000');
		setFontSize(parseFloat(localStorage.getItem('subtitles-font-size') || 4));
		setFontShadowWidth(parseInt(localStorage.getItem('subtitles-font-shadow-width') || 2));
	},

	setFontColor: (fontColor) => {
		set({ fontColor });
		localStorage.setItem('subtitles-font-color', fontColor);
	},

	setFontShadowColor: (fontShadowColor) => {
		set({ fontShadowColor });
		localStorage.setItem('subtitles-font-shadow-color', fontShadowColor);
	},

	setFontSize: (fontSize) => {
		set({ fontSize });
		localStorage.setItem('subtitles-font-size', fontSize);
	},

	setFontShadowWidth: (fontShadowWidth) => {
		set({ fontShadowWidth });
		localStorage.setItem('subtitles-font-shadow-width', fontShadowWidth);
	},

	setSelectedSubtitleUrl: (selectedSubtitleUrl) => set({ selectedSubtitleUrl }),
	setSubtitleOffsetMillis: (subtitleOffsetMillis) => set({ subtitleOffsetMillis }),

	toggleSubtitlesControls: () => {
		let { controlsShown } = get();
		set({ controlsShown: !controlsShown });
	},

	hideSubtitlesControls: () => {
		set({ controlsShown: false });
	},

	closeAll: () => {
		set({ controlsShown: false });
	},
}));
