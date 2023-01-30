export default class Client {
	static baseUrl = 'http://localhost:8080/api';

	static getTags = async () => {
		return await fetch(`${Client.baseUrl}/tags`).then((response) => response.json());
	};

	static getTag = async ({ tagId }) => {
		return await fetch(`${Client.baseUrl}/tags/${tagId}`).then((response) => response.json());
	};

	static saveTag = async (tag, successCallback) => {
		return await fetch(`${Client.baseUrl}/tags/${tag.id}`, {
			method: 'POST',
			body: JSON.stringify(tag),
		});
	};

	static getItems = async () => {
		return await fetch(`${Client.baseUrl}/items`).then((response) => response.json());
	};

	static getItem(itemId, successCallback) {
		fetch(`${Client.baseUrl}/items/${itemId}`)
			.then((response) => response.json())
			.then((item) => successCallback(item));
	}

	static addAnnotationToTag = async ({ tagId, annotation }) => {
		return await fetch(`${Client.baseUrl}/tags/${tagId}/annotations`, {
			method: 'POST',
			body: JSON.stringify(annotation),
		});
	};

	static removeAnnotationFromTag = async ({ tagId, annotationId }) => {
		return await fetch(`${Client.baseUrl}/tags/${tagId}/annotations/${annotationId}`, {
			method: 'DELETE',
		});
	};

	static getAvailableAnnotations = async (tagId) => {
		return await fetch(`${Client.baseUrl}/tags/${tagId}/available-annotations`).then((response) => response.json());
	};

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

	static refreshCovers() {
		fetch(`${Client.baseUrl}/items/refresh-covers`);
	}

	static refreshPreview() {
		fetch(`${Client.baseUrl}/items/refresh-preview`);
	}

	static getExportMetadataUrl() {
		return `${Client.baseUrl}/export-metadata.json`;
	}

	static buildFileUrl(storagePath) {
		return `${Client.baseUrl}/file/${encodeURIComponent(storagePath)}`;
	}
}
