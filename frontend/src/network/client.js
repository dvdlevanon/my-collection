export default class Client {
	static baseUrl = 'http://localhost:8080';

	static getTags(successCallback) {
		fetch(`${Client.baseUrl}/tags`)
			.then((response) => response.json())
			.then((tags) => successCallback(tags));
	}

	static saveTag(tag, successCallback) {
		fetch(`${Client.baseUrl}/tags/${tag.id}`, {
			method: 'POST',
			body: JSON.stringify(tag),
		}).then(successCallback);
	}

	static getItems(successCallback) {
		fetch(`${Client.baseUrl}/items`)
			.then((response) => response.json())
			.then((items) => successCallback(items));
	}

	static getItem(itemId, successCallback) {
		fetch(`${Client.baseUrl}/items/${itemId}`)
			.then((response) => response.json())
			.then((item) => successCallback(item));
	}

	static saveItem(item, successCallback) {
		fetch(`${Client.baseUrl}/items/${item.id}`, {
			method: 'POST',
			body: JSON.stringify(item),
		}).then(successCallback);
	}

	static removeTagFromItem(itemId, tagId, successCallback) {
		fetch(`${Client.baseUrl}/items/${itemId}/remove-tag/${tagId}`, {
			method: 'POST',
		}).then(successCallback);
	}

	static uploadFile(storagePath, file, successCallback) {
		let formData = new FormData();
		formData.append('path', storagePath);
		formData.append('file', file);

		fetch(`${Client.baseUrl}/upload-file`, {
			method: 'POST',
			body: formData,
		})
			.then((response) => response.json())
			.then((fileUrl) => successCallback(fileUrl));
	}

	static refreshPreview() {
		fetch(`${Client.baseUrl}/items/refresh-preview`);
	}

	static buildStorageUrl(storagePath) {
		return `${Client.baseUrl}/storage/${encodeURIComponent(storagePath)}`;
	}

	static buildStreamUrl(streamPath) {
		return `${Client.baseUrl}/stream/${encodeURIComponent(streamPath)}`;
	}
}
