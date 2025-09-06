import { create } from 'zustand';

export const usePlayerActionStore = create((set, get) => ({
	cropMode: false,
	cropFrame: null,
	startHighlightSecond: -1,

	startHighlightCreation: (second) => {
		set({ startHighlightSecond: second });
	},

	highlightActive: () => {
		const { startHighlightSecond } = get();
		return startHighlightSecond !== -1;
	},

	highlightCompleted: () => {
		const { startHighlightSecond } = get();
		set({ startHighlightSecond: -1 });
		return startHighlightSecond;
	},

	highlightCanceled: () => {
		set({ startHighlightSecond: -1 });
	},

	startCrop: () => {
		set({ cropMode: true });
	},

	cropActive: () => {
		const { cropMode } = get();
		return cropMode;
	},

	cropCompleted: () => {
		const { cropFrame } = get();
		set({ cropMode: false });
		return cropFrame;
	},

	cropCanceled: () => {
		set({ cropMode: false });
	},

	setCropFrame: (cropFrame) => set({ cropFrame }),

	closeAll: () => {
		set({ cropMode: false });
		set({ startHighlightSecond: -1 });
	},
}));
