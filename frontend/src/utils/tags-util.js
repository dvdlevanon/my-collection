import Client from './client';

export default class TagsUtil {
	static directoriesTag;
	static dailyMixTag;
	static highlightsTag;
	static specTag;
	static categories = [];

	static initCategories(categories) {
		TagsUtil.categories = categories;
	}

	static initSpecialTags(specialTags) {
		for (let i = 0; i < specialTags.length; i++) {
			if (specialTags[i].title === 'Directories') {
				// directories.go
				TagsUtil.directoriesTag = specialTags[i];
			} else if (specialTags[i].title === 'DailyMix') {
				// automix.go
				TagsUtil.dailyMixTag = specialTags[i];
			} else if (specialTags[i].title === 'Spec') {
				// spectagger.go
				TagsUtil.specTag = specialTags[i];
			} else if (specialTags[i].title === 'Highlights') {
				// highlights.go
				TagsUtil.highlightsTag = specialTags[i];
			}
		}

		if (!TagsUtil.directoriesTag || !TagsUtil.dailyMixTag || !TagsUtil.highlightsTag || !TagsUtil.specTag) {
			console.log('Missing mandatory special tags');
		}
	}

	static isDirectoriesCategory(tagId) {
		return tagId === TagsUtil.directoriesTag.id;
	}

	static isDailymixCategory(tagId) {
		return tagId === TagsUtil.dailyMixTag.id;
	}

	static isHighlightsCategory(tagId) {
		return tagId === TagsUtil.highlightsTag.id;
	}

	static isSpecCategory(tagId) {
		return tagId === TagsUtil.specTag.id;
	}

	static isSpecialCategory(tagId) {
		return (
			TagsUtil.isDirectoriesCategory(tagId) ||
			TagsUtil.isDailymixCategory(tagId) ||
			TagsUtil.isSpecCategory(tagId) ||
			TagsUtil.isHighlightsCategory(tagId)
		);
	}

	static showAsThumbnail(tagId) {
		if (!TagsUtil.categories) {
			return false;
		}

		let category = TagsUtil.categories.find((cur) => cur.id == tagId);
		return category.display_style === 'portrait';
	}

	static showAsBanner(tagId) {
		if (!TagsUtil.categories) {
			return false;
		}

		let category = TagsUtil.categories.find((cur) => cur.id == tagId);
		return category.display_style === 'banner';
	}

	static getBannerCategoryId() {
		if (!TagsUtil.categories) {
			return 0;
		}

		return TagsUtil.categories.find((cur) => cur.display_style === 'banner').id;
	}

	static allowToAddToCategory(tagId) {
		return !(
			TagsUtil.isDirectoriesCategory(tagId) ||
			TagsUtil.isDailymixCategory(tagId) ||
			TagsUtil.isSpecCategory(tagId)
		);
	}

	static allowToSetImageToCategory(tagId) {
		return !(
			TagsUtil.isDirectoriesCategory(tagId) ||
			TagsUtil.isDailymixCategory(tagId) ||
			TagsUtil.isSpecCategory(tagId)
		);
	}

	static getCategories(tags) {
		if (!tags) {
			return [];
		}

		let result = tags.filter((tag) => {
			return !tag.parentId;
		});

		return result;
	}

	static tagTitleToFileName(title) {
		return title.toLowerCase().replaceAll(' ', '-');
	}

	static normalizeTagTitle(rawTitle) {
		let regex = /\b[A-Z]{2,}\b/g;
		let noConsequensiveUpperCaser = rawTitle.replace(regex, function (match) {
			return match.toLowerCase();
		});

		let title = noConsequensiveUpperCaser
			.replaceAll('-', ' ')
			.replaceAll('_', ' ')
			.replaceAll('.', ' ')
			.replaceAll(',', ' ')
			.replaceAll('[', '')
			.replaceAll(']', '')
			.replaceAll('(', '')
			.replaceAll(')', '')
			.replaceAll(/([A-Z])/g, ' $1')
			.trim();

		while (title.includes('  ')) {
			title = title.replaceAll('  ', ' ');
		}

		if (!title) {
			return '';
		}

		const words = title.split(' ');

		for (let i = 0; i < words.length; i++) {
			words[i] = words[i][0].toUpperCase() + words[i].substr(1);
		}

		return words.join(' ');
	}

	static hasImage = (tag) => {
		if (this.isSpecialCategory(tag.parentId)) {
			return true;
		}

		return tag.images && tag.images.length > 0;
	};

	static hasTagImage = (tag, tit) => {
		if (TagsUtil.isDirectoriesCategory(tag.parentId)) {
			return true;
		} else if (TagsUtil.isDailymixCategory(tag.parentId)) {
			return true;
		}

		if (tit && tag.images) {
			let selectedImage = tag.images.find((image) => image.imageType === tit.id);
			if (selectedImage && selectedImage.url) {
				return true;
			}
		}

		return false;
	};

	static getTagImageThumbnailRect = (tag, selectedTit) => {
		if (selectedTit && tag.images) {
			let selectedImage = tag.images.find((image) => image.imageType === selectedTit.id);
			if (selectedImage) {
				return selectedImage.thumbnail_rect;
			}
		}

		return {
			x: 0,
			y: 0,
			width: 200,
			height: 200,
		};
	};

	static getNoBannerImageUrl = () => {
		return Client.buildFileUrl(Client.buildInternalStoragePath('tags-image-types/banner/1.jpg'));
	};

	static getTagImageUrl = (tag, selectedTit, noFallback) => {
		if (TagsUtil.isDirectoriesCategory(tag.parentId)) {
			return Client.buildFileUrl(Client.buildInternalStoragePath('tags-image-types/directory/directory.png'));
		} else if (TagsUtil.isDailymixCategory(tag.parentId)) {
			return Client.buildFileUrl(Client.buildInternalStoragePath('tags-image-types/dailymix/dailymix.png'));
		}

		if (selectedTit && tag.images) {
			let selectedImage = tag.images.find((image) => image.imageType === selectedTit.id);
			if (selectedImage && selectedImage.url) {
				return Client.buildFileUrl(selectedImage.url);
			}
		}

		if (!noFallback && tag.images) {
			for (let i = 0; i < tag.images.length; i++) {
				if (tag.images[i].url) {
					return Client.buildFileUrl(tag.images[i].url);
				}
			}
		}

		return Client.buildFileUrl(Client.buildInternalStoragePath('tags-image-types/none/1.jpg'));
	};

	static sortByTitle = (tags) => {
		return tags.sort((a, b) => (a.title > b.title ? 1 : a.title < b.title ? -1 : 0));
	};

	static itemsCount = (tag) => {
		return tag.items ? tag.items.length : 0;
	};
}
