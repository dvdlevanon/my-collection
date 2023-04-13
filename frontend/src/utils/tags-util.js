export default class TagsUtil {
	static isDirectoriesCategory(tagId) {
		return tagId == 1; // directories.go
	}

	static isDailymixCategory(tagId) {
		return tagId == 343; // automix.go
	}

	static isSpecialCategory(tagId) {
		return TagsUtil.isDirectoriesCategory(tagId) || TagsUtil.isDailymixCategory(tagId);
	}

	static getCategories(tags) {
		if (!tags) {
			return [];
		}

		return tags.filter((tag) => {
			return !tag.parentId;
		});
	}

	static normalizeTagTitle(rawTitle) {
		let title = rawTitle
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

		const words = title.split(' ');
		for (let i = 0; i < words.length; i++) {
			words[i] = words[i][0].toUpperCase() + words[i].substr(1);
		}

		return words.join(' ');
	}
}
