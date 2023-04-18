export default class TagsUtil {
	static directoriesTag;
	static dailyMixTag;

	static initSpecialTags(specialTags) {
		for (let i = 0; i < specialTags.length; i++) {
			if (specialTags[i].title == 'Directories') {
				// directories.go
				TagsUtil.directoriesTag = specialTags[i];
			} else if (specialTags[i].title == 'DailyMix') {
				// automix.go
				TagsUtil.dailyMixTag = specialTags[i];
			}
		}

		if (!TagsUtil.directoriesTag || !TagsUtil.dailyMixTag) {
			console.log('Missing mandatory special tags');
		}
	}

	static isDirectoriesCategory(tagId) {
		return tagId == TagsUtil.directoriesTag.id;
	}

	static isDailymixCategory(tagId) {
		return tagId == TagsUtil.dailyMixTag.id;
	}

	static isSpecialCategory(tagId) {
		return TagsUtil.isDirectoriesCategory(tagId) || TagsUtil.isDailymixCategory(tagId);
	}

	static getCategories(tags) {
		if (!tags) {
			return [];
		}

		let result = tags.filter((tag) => {
			return !tag.parentId;
		});

		console.log(result);
		return result;
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
