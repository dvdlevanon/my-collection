export default class Client {
	static baseUrl = 'http://localhost:8080';
	static apiUrl = `${Client.baseUrl}/api`;

	static getTags = async () => {
		return await fetch(`${Client.apiUrl}/tags`).then((response) => response.json());
	};

	static getTag = async ({ tagId }) => {
		return await fetch(`${Client.apiUrl}/tags/${tagId}`).then((response) => response.json());
	};

	static createTag = async (tag) => {
		return await fetch(`${Client.apiUrl}/tags`, {
			method: 'POST',
			body: JSON.stringify(tag),
		});
	};

	static saveTag = async (tag) => {
		return await fetch(`${Client.apiUrl}/tags/${tag.id}`, {
			method: 'POST',
			body: JSON.stringify(tag),
		});
	};

	static removeTag = async (tagId) => {
		return await fetch(`${Client.apiUrl}/tags/${tagId}`, {
			method: 'DELETE',
		});
	};

	static getItems = async () => {
		return await fetch(`${Client.apiUrl}/items`).then((response) => response.json());
	};

	static getItem = async (itemId) => {
		return await fetch(`${Client.apiUrl}/items/${itemId}`).then((response) => response.json());
	};

	static addAnnotationToTag = async ({ tagId, annotation }) => {
		return await fetch(`${Client.apiUrl}/tags/${tagId}/annotations`, {
			method: 'POST',
			body: JSON.stringify(annotation),
		});
	};

	static removeAnnotationFromTag = async ({ tagId, annotationId }) => {
		return await fetch(`${Client.apiUrl}/tags/${tagId}/annotations/${annotationId}`, {
			method: 'DELETE',
		});
	};

	static getAvailableAnnotations = async (tagId) => {
		return await fetch(`${Client.apiUrl}/tags/${tagId}/available-annotations`).then((response) => response.json());
	};

	static saveItem(item, successCallback) {
		fetch(`${Client.apiUrl}/items/${item.id}`, {
			method: 'POST',
			body: JSON.stringify(item),
		}).then(successCallback);
	}

	static removeTagFromItem(itemId, tagId, successCallback) {
		fetch(`${Client.apiUrl}/items/${itemId}/remove-tag/${tagId}`, {
			method: 'POST',
		}).then(successCallback);
	}

	static uploadFile(storagePath, file, successCallback) {
		let formData = new FormData();
		formData.append('path', storagePath);
		formData.append('file', file);

		fetch(`${Client.apiUrl}/upload-file`, {
			method: 'POST',
			body: formData,
		})
			.then((response) => response.json())
			.then((fileUrl) => successCallback(fileUrl));
	}

	static refreshCovers() {
		fetch(`${Client.apiUrl}/items/refresh-covers`);
	}

	static refreshPreview() {
		fetch(`${Client.apiUrl}/items/refresh-preview`);
	}

	static refreshVideoMetadata() {
		fetch(`${Client.apiUrl}/items/refresh-video-metadata`);
	}

	static getExportMetadataUrl() {
		return `${Client.apiUrl}/export-metadata.json`;
	}

	static buildFileUrl(storagePath) {
		return `${Client.apiUrl}/file/${encodeURIComponent(storagePath)}`;
	}

	static buildInternalStoragePath(storagePath) {
		return `.internal-storage/${storagePath}`;
	}
}
