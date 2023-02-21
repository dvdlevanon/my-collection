export default class Client {
	static baseUrl = 'http://localhost:8080';
	static apiUrl = `${Client.baseUrl}/api`;

	static getTags = async () => {
		return await fetch(`${Client.apiUrl}/tags`).then((response) => response.json());
	};

	static getDirectories = async () => {
		return await fetch(`${Client.apiUrl}/directories`).then((response) => response.json());
	};

	static addOrUpdateDirectory = async (directory) => {
		return await fetch(`${Client.apiUrl}/directories`, {
			method: 'POST',
			body: JSON.stringify(directory),
		});
	};

	static setDirectoryCategories = async (directory) => {
		return await fetch(`${Client.apiUrl}/directories/tags/${directory.path}`, {
			method: 'POST',
			body: JSON.stringify(directory),
		});
	};

	static removeDirectory = async (directoryPath) => {
		return await fetch(`${Client.apiUrl}/directories/${directoryPath}`, {
			method: 'DELETE',
		});
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

	static imageDirectoryChoosen = async (tagId, directoryPath) => {
		return await fetch(`${Client.apiUrl}/tags/${tagId}/auto-image`, {
			method: 'POST',
			body: JSON.stringify({ url: directoryPath }),
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

	static setMainCover = async (itemId, second) => {
		return await fetch(`${Client.apiUrl}/items/${itemId}/main-cover?second=${second}`, {
			method: 'POST',
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
