import queryString from 'query-string';

export default class GalleryUrlParams {
	constructor(searchParams, setSearchParams) {
		this.searchParams = searchParams;
		this.setSearchParams = setSearchParams;
	}

	static buildUrlParams(activeTagId) {
		let params = [];
		params['active-tags'] = activeTagId + '=true';
		return queryString.stringify(params);
	}

	getUrlParamsString() {
		return this.searchParams.toString();
	}

	getTags() {
		let parsed = queryString.parse(this.searchParams.toString());
		let activeTags = parsed['active-tags'];

		if (!activeTags) {
			return [];
		}

		return activeTags.split(',').map((tag) => {
			let parts = tag.split('=');
			let tagId = parseInt(parts[0]);
			let selected = parts[1] == 'true';
			return { id: tagId, selected: selected };
		});
	}

	getActiveTags() {
		let tagsStatus = this.getTags();
		return tagsStatus.map((tag) => tag.id);
	}

	getSelectedTags() {
		let tagsStatus = this.getTags();
		return tagsStatus.filter((tag) => tag.selected).map((tag) => tag.id);
	}

	tagToString(tagId, selected) {
		return tagId + '=' + selected;
	}

	toggleTagSelection(tagId) {
		let tagsStatus = this.getTags();
		let updatedTags = tagsStatus.map((tag) => {
			return this.tagToString(tag.id, tagId == tag.id ? !tag.selected : tag.selected);
		});
		this.updateActiveTagsString(updatedTags);
	}

	activateTag(tagId, keepSelection) {
		let tagsStatus = this.getTags();
		let updatedTags = tagsStatus.map((tag) => this.tagToString(tag.id, keepSelection ? tag.selected : false));
		updatedTags.push(this.tagToString(tagId, true));
		this.updateActiveTagsString(updatedTags);
	}

	buildActivateTagUrl(tagId) {
		let tagsStatus = this.getTags();
		let updatedTags = tagsStatus.map((tag) => this.tagToString(tag.id, false));
		updatedTags.push(this.tagToString(tagId, true));
		return this.buildActiveTagsString(updatedTags);
	}

	deactivateTag(tagId) {
		let tagsStatus = this.getTags();
		let updatedTags = tagsStatus
			.filter((tag) => tag.id != tagId)
			.map((tag) => this.tagToString(tag.id, tag.selected));
		this.updateActiveTagsString(updatedTags);
	}

	deactivateAllTags() {
		this.updateActiveTagsString([]);
	}

	getConditionType() {
		let parsed = queryString.parse(this.searchParams.toString());
		return parsed['condition-type'] || '||';
	}

	setConditionType(conditionType) {
		let parsed = queryString.parse(this.searchParams.toString());
		parsed['condition-type'] = conditionType;
		this.setSearchParams(parsed);
	}

	updateActiveTagsString(activeTags) {
		this.setSearchParams(this.buildActiveTagsString(activeTags));
	}

	buildActiveTagsString(activeTags) {
		let parsed = queryString.parse(this.searchParams.toString());
		if (activeTags.length > 0) {
			parsed['active-tags'] = activeTags.join(',');
		} else {
			delete parsed['active-tags'];
		}

		return queryString.stringify(parsed);
	}
}
