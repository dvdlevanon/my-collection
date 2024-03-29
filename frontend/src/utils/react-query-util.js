export default class ReactQueryUtil {
	static SPECIAL_TAGS_KEY = ['special-tags'];
	static CATEGORIES_KEY = ['categories'];
	static TAGS_KEY = ['tags'];
	static ITEMS_KEY = ['items'];
	static QUEUE_METADATA_KEY = ['queue-metadata'];
	static STATS_KEY = ['stats'];
	static DIRECTORIES_KEY = ['directories'];
	static TAG_IMAGE_TYPES_KEY = ['tag-image-types'];

	static availableAnnotationsKey = (tagId) => {
		return ['available-annotations', { id: tagId }];
	};

	static tagCustomCommands = (tagId) => {
		return ['tag-custom-commands', { id: tagId }];
	};

	static itemKey = (itemId) => {
		return ['items', { id: itemId }];
	};

	static tagKey = (tagId) => {
		return ['tags', { id: tagId }];
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
