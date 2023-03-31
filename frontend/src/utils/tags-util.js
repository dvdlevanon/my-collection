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
}
