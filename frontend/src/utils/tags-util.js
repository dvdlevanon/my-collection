import Client from './client';

export default class TagsUtil {
	static directoriesTag;
	static dailyMixTag;
	static highlightsTag;
	static specTag;

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

	static getTagImageUrl = (tag, selectedTit) => {
		if (TagsUtil.isDirectoriesCategory(tag.parentId)) {
			return Client.buildFileUrl(Client.buildInternalStoragePath('tags-image/directory/directory.png'));
		} else if (TagsUtil.isDailymixCategory(tag.parentId)) {
			return Client.buildFileUrl(Client.buildInternalStoragePath('tags-image/dailymix/dailymix.png'));
		}

		if (selectedTit && tag.images) {
			let selectedImage = tag.images.find((image) => image.imageType === selectedTit.id);
			if (selectedImage && selectedImage.url) {
				return Client.buildFileUrl(selectedImage.url);
			}
		}

		if (tag.images) {
			for (let i = 0; i < tag.images.length; i++) {
				if (tag.images[i].url) {
					return Client.buildFileUrl(tag.images[i].url);
				}
			}
		}

		return Client.buildFileUrl(Client.buildInternalStoragePath('tags-image/none/1.jpg'));
	};

	static sortByTitle = (tags) => {
		return tags.sort((a, b) => (a.title > b.title ? 1 : a.title < b.title ? -1 : 0));
	};

	static itemsCount = (tag) => {
		return tag.items ? tag.items.length : 0;
	};
}
