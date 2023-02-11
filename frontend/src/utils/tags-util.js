export default class TagsUtil {
	static isDirectoriesCategory(tagId) {
		return tagId == 1; // directories.go
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
