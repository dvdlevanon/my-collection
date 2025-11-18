import { produce } from 'immer';
import Client from './client';

export default class ReactQueryUtil {
	static ASYNC_TASKS_TIMEOUT = 1000;
	static SPECIAL_TAGS_KEY = ['special-tags'];
	static CATEGORIES_KEY = ['categories'];
	static TAGS_KEY = ['tags'];
	static ITEMS_KEY = ['items'];
	static QUEUE_METADATA_KEY = ['queue-metadata'];
	static STATS_KEY = ['stats'];
	static DIRECTORIES_KEY = ['directories'];
	static TAG_IMAGE_TYPES_KEY = ['tag-image-types'];

	static availableAnnotationsKey = (tagId) => {
		return ['available-annotations', { id: String(tagId) }];
	};

	static tagCustomCommands = (tagId) => {
		return ['tag-custom-commands', { id: String(tagId) }];
	};

	static itemKey = (itemId) => {
		return ['items', { id: String(itemId) }];
	};

	static subtitleQueryKey = (url) => {
		return ['subtitle', { id: url }];
	};

	static availableSubtitleQueryKey = (itemId) => {
		return ['available-subtitle', { id: String(itemId) }];
	};

	static onlineSubtitleQueryKey = (itemId, lang, aiTranslated) => {
		return ['online-subtitle', { id: String(itemId), lang: lang, aiTranslated: String(aiTranslated) }];
	};

	static tagKey = (tagId) => {
		return ['tags', { id: String(tagId) }];
	};

	static suggestedItemsKey = (itemId) => {
		return ['suggested', { id: String(itemId) }];
	};

	static tasksPageKey = (pageId, pageSize) => {
		return ['tasks', { page: pageId, pageSize: pageSize }];
	};

	static updateItem = (queryClient, itemId, withDelay) => {
		if (withDelay) {
			setTimeout(() => {
				queryClient.invalidateQueries({ queryKey: ReactQueryUtil.itemKey(itemId) });
			}, ReactQueryUtil.ASYNC_TASKS_TIMEOUT);
		} else {
			queryClient.invalidateQueries({ queryKey: ReactQueryUtil.itemKey(itemId) });
		}
	};

	static updateTags = (queryClient, tagId, updater) => {
		queryClient.setQueryData(ReactQueryUtil.TAGS_KEY, (oldTags) =>
			produce(oldTags, (draft) => {
				const tagToUpdate = draft.find((tag) => tag.id === tagId);
				if (tagToUpdate) {
					updater(tagToUpdate);
				}
			})
		);
	};

	static getTag = (queryClient, tagId) => {
		const tags = queryClient.getQueryData(ReactQueryUtil.TAGS_KEY);
		return tags.find((tag) => tag.id === tagId);
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

	static tagsQuery = () => {
		return {
			queryKey: ReactQueryUtil.TAGS_KEY,
			queryFn: Client.getTags,
		};
	};

	static itemQuery = (itemId) => {
		return {
			queryKey: ReactQueryUtil.itemKey(itemId),
			queryFn: () => Client.getItem(itemId),
		};
	};

	static suggestionQuery = (itemId) => {
		return {
			queryKey: ReactQueryUtil.suggestedItemsKey(itemId),
			queryFn: () => Client.getSuggestedItems(itemId),
			staleTime: Infinity,
			cacheTime: Infinity,
		};
	};

	static availableAnnotationsQuery = (categoryId) => {
		return {
			queryKey: ReactQueryUtil.availableAnnotationsKey(categoryId),
			queryFn: () => Client.getAvailableAnnotations(categoryId),
		};
	};

	static subtitleQuery = (url) => {
		return {
			queryKey: ReactQueryUtil.subtitleQueryKey(url),
			queryFn: () => Client.getSubtitle(url),
		};
	};

	static availableSubtitleQuery = (itemId) => {
		return {
			queryKey: ReactQueryUtil.availableSubtitleQueryKey(itemId),
			queryFn: () => Client.getAvailableSubtitle(itemId),
		};
	};

	static onlineSubtitleQuery = (itemId, lang, aiTranslated) => {
		return {
			queryKey: ReactQueryUtil.onlineSubtitleQueryKey(itemId, lang, aiTranslated),
			queryFn: () => Client.getOnlineSubtitle(itemId, lang, aiTranslated),
		};
	};
}
