export default class ReactQueryUtil {
	static TAGS_KEY = ['tags'];
	static ITEMS_KEY = ['items'];
	static QUEUE_METADATA_KEY = ['queue-metadata'];
	static DIRECTORIES_KEY = ['directories'];

	static availableAnnotationsKey = (tagId) => {
		return ['available-annotations', { id: tagId }];
	};

	static itemKey = (itemId) => {
		return ['items', { id: itemId }];
	};

	static suggestedItemsKey = (itemId) => {
		return ['suggested', { id: itemId }];
	};

	static tasksPageKey = (pageId, pageSize) => {
		return ['tasks', { page: pageId, pageSize: pageSize }];
	};

	static updateTags = (queryClient, tagId, updater) => {
		queryClient.setQueryData(ReactQueryUtil.TAGS_KEY, (oldTags) => {
			return oldTags.map((cur) => {
				if (cur.id != tagId) {
					return cur;
				}

				return updater(cur);
			});
		});
	};

	static updateDirectories = (queryClient, directoryPath, updater) => {
		queryClient.setQueryData(ReactQueryUtil.DIRECTORIES_KEY, (oldDirectories) => {
			return oldDirectories.map((cur) => {
				if (cur.path != directoryPath) {
					return cur;
				}

				return updater(cur);
			});
		});
	};
}
