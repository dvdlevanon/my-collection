export default class ReactQueryUtil {
	static TAGS_KEY = ['tags'];
	static ITEMS_KEY = ['items'];

	static availableAnnotationsKey = (tagId) => {
		return ['available-annotations', { id: tagId }];
	};

	static itemKey = (itemId) => {
		return ['items', { id: itemId }];
	};

	static updateTags = (queryClient, tagId, updater) => {
		queryClient.setQueryData(['tags'], (oldTags) => {
			return oldTags.map((cur) => {
				if (cur.id != tagId) {
					return cur;
				}

				return updater(cur);
			});
		});
	};
}
