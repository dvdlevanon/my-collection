import { create } from 'zustand';

const sortAnnotations = (annotations) => {
	return annotations.sort((a, b) => (a.title > b.title ? 1 : a.title < b.title ? -1 : 0));
};

export const useTagAnnotationsStore = create((set, get) => ({
	tag: null,
	availableAnnotations: [],

	setTag: (tag) => set({ tag }),
	setAvailableAnnotations: (availableAnnotations) => set({ availableAnnotations }),

	getAttachedAnnotations: () => {
		const { tag, availableAnnotations } = get();

		if (!availableAnnotations || !tag.tags_annotations) {
			return [];
		}

		let attachedAnnotations = availableAnnotations.filter((annotation) => {
			return tag.tags_annotations.some((cur) => cur.id == annotation.id);
		});

		return sortAnnotations(attachedAnnotations);
	},

	getUnattachedAnnotations: () => {
		const { tag, availableAnnotations } = get();

		if (!availableAnnotations) {
			return [];
		}

		let unattachedAnnotations = availableAnnotations.filter((annotation) => {
			if (!tag.tags_annotations) return true;
			return !tag.tags_annotations.some((cur) => cur.id == annotation.id);
		});

		return sortAnnotations(unattachedAnnotations);
	},
}));
