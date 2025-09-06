export default class Client {
	static baseUrl = window.location.hostname + (window.location.port ? ':' + window.location.port : '');
	static apiUrl = process.env.REACT_APP_API_URL || `http://${Client.baseUrl}/api`;
	static websocketUrl = process.env.REACT_APP_WEBSOCKET_URL || `ws://${Client.baseUrl}/api/ws`;

	static getSpecialTags = async () => {
		return await fetch(`${Client.apiUrl}/special-tags`).then((response) => response.json());
	};

	static getCategories = async () => {
		return await fetch(`${Client.apiUrl}/categories`).then((response) => response.json());
	};

	static getTags = async () => {
		return await fetch(`${Client.apiUrl}/tags`).then((response) => response.json());
	};

	static getQueueMetadata = async () => {
		return await fetch(`${Client.apiUrl}/queue/metadata`).then((response) => response.json());
	};

	static getStats = async () => {
		return await fetch(`${Client.apiUrl}/stats`).then((response) => response.json());
	};

	static getTasks = async (page, pageSize) => {
		return await fetch(`${Client.apiUrl}/queue/tasks?page=${page}&pageSize=${pageSize}`).then((response) =>
			response.json()
		);
	};

	static clearFinishedTasks = async () => {
		return await fetch(`${Client.apiUrl}/queue/clear-finished`, {
			method: 'POST',
		});
	};

	static continueProcessingTasks = async () => {
		return await fetch(`${Client.apiUrl}/queue/continue`, {
			method: 'POST',
		});
	};

	static pauseProcessingTasks = async () => {
		return await fetch(`${Client.apiUrl}/queue/pause`, {
			method: 'POST',
		});
	};

	static setDirectoryCategories = async (directory) => {
		return await fetch(`${Client.apiUrl}/directories/tags/${directory.path}`, {
			method: 'POST',
			body: JSON.stringify(directory),
		});
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

	static excludeFromRandomMix = async (tag) => {
		return await fetch(`${Client.apiUrl}/tags/${tag.id}/random-mix/exclude`, {
			method: 'POST',
		});
	};

	static includeInRandomMix = async (tag) => {
		return await fetch(`${Client.apiUrl}/tags/${tag.id}/random-mix/include`, {
			method: 'POST',
		});
	};

	static removeTag = async (tagId) => {
		return await fetch(`${Client.apiUrl}/tags/${tagId}`, {
			method: 'DELETE',
		});
	};

	static mixOnDemand = async (desc, tags) => {
		return await fetch(`${Client.apiUrl}/mix-on-demand?desc=${encodeURIComponent(desc)}`, {
			method: 'POST',
			body: JSON.stringify(tags),
		}).then((response) => response.json());
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

	static getSubtitle = async (itemId, subtitleName) => {
		return await fetch(`${Client.apiUrl}/subtitles/${itemId}?name=${encodeURIComponent(subtitleName)}`).then(
			(response) => response.json()
		);
	};

	static getAvailableSubtitle = async (itemId) => {
		return await fetch(`${Client.apiUrl}/subtitles/${itemId}/available`).then((response) => response.json());
	};

	static getSuggestedItems = async (itemId) => {
		return await fetch(`${Client.apiUrl}/items/${itemId}/suggestions`).then((response) => response.json());
	};

	static getItem = async (itemId) => {
		return await fetch(`${Client.apiUrl}/items/${itemId}`).then((response) => response.json());
	};

	static getTag = async (tagId) => {
		return await fetch(`${Client.apiUrl}/tags/${tagId}`).then((response) => response.json());
	};

	static getItemLocation = async (itemId) => {
		return await fetch(`${Client.apiUrl}/items/${itemId}/location`).then((response) => response.json());
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

	static deleteItem = async (itemId, deleteRealFile) => {
		return await fetch(`${Client.apiUrl}/items/${itemId}?deleteRealFile=${deleteRealFile}`, {
			method: 'DELETE',
		});
	};

	static forceProcessItem = async (itemId) => {
		return await fetch(`${Client.apiUrl}/items/${itemId}/process`, {
			method: 'POST',
		});
	};

	static optimizeItem = async (itemId) => {
		return await fetch(`${Client.apiUrl}/items/${itemId}/optimize`, {
			method: 'POST',
		});
	};

	static splitItem = async (itemId, second) => {
		return await fetch(`${Client.apiUrl}/items/${itemId}/split?second=${second}`, {
			method: 'POST',
		});
	};

	static makeHighlight = async (itemId, startSecond, endSecond, highlightId) => {
		return await fetch(
			`${Client.apiUrl}/items/${itemId}/make-highlight?start=${startSecond}&end=${endSecond}&highlight-id=${highlightId}`,
			{
				method: 'POST',
			}
		);
	};

	static cropFrame = async (itemId, second, crop) => {
		return await fetch(
			`${Client.apiUrl}/items/${itemId}/crop-frame?second=${second}&crop-x=${crop.x}&crop-y=${crop.y}&crop-width=${crop.width}&crop-height=${crop.height}`,
			{
				method: 'POST',
			}
		);
	};

	static removeTagImageFromTag = async (tagId, titId) => {
		return await fetch(`${Client.apiUrl}/tags/${tagId}/tit/${titId}`, {
			method: 'DELETE',
		});
	};

	static removeAnnotationFromTag = async ({ tagId, annotationId }) => {
		return await fetch(`${Client.apiUrl}/tags/${tagId}/annotations/${annotationId}`, {
			method: 'DELETE',
		});
	};

	static updateTagImage = async (tagId, image) => {
		return await fetch(`${Client.apiUrl}/tags/${tagId}/images/${image.id}`, {
			method: 'POST',
			body: JSON.stringify(image),
		});
	};

	static getAvailableAnnotations = async (tagId) => {
		return await fetch(`${Client.apiUrl}/tags/${tagId}/available-annotations`).then((response) => response.json());
	};

	static getTagCustomCommands = async (tagId) => {
		return await fetch(`${Client.apiUrl}/tags/${tagId}/tag-custom-commands`).then((response) => response.json());
	};

	static getTagImageTypes = async () => {
		return await fetch(`${Client.apiUrl}/tags/tag-image-types`).then((response) => response.json());
	};

	static getFsDir = async (path, depth) => {
		return await fetch(`${Client.apiUrl}/fs?path=${encodeURIComponent(path)}&depth=${depth}`).then((response) =>
			response.json()
		);
	};

	static includeDir = async (path, subdirs, hierarchy) => {
		return fetch(
			`${Client.apiUrl}/fs/include?path=${encodeURIComponent(path)}&subdirs=${subdirs}&hierarchy=${hierarchy}`,
			{
				method: 'POST',
			}
		);
	};

	static excludeDir = async (path) => {
		return fetch(`${Client.apiUrl}/fs/exclude?path=${encodeURIComponent(path)}`, {
			method: 'POST',
		});
	};

	static fetchTextFile = async (path) => {
		return await fetch(path).then((response) => response.text());
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

	static uploadFileFromUrl(storagePath, fileUrl, successCallback) {
		fetch(
			`${Client.apiUrl}/upload-file-from-url?url=${encodeURIComponent(fileUrl)}&path=${encodeURIComponent(
				storagePath
			)}`,
			{
				method: 'POST',
			}
		)
			.then((response) => response.json())
			.then((fileUrl) => successCallback(fileUrl));
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

	static refreshCovers(force) {
		fetch(`${Client.apiUrl}/items/refresh-covers?force=${force}`, { method: 'POST' });
	}

	static refreshPreview(force) {
		fetch(`${Client.apiUrl}/items/refresh-preview?force=${force}`, { method: 'POST' });
	}

	static refreshVideoMetadata(force) {
		fetch(`${Client.apiUrl}/items/refresh-video-metadata?force=${force}`, { method: 'POST' });
	}

	static refreshFileMetadata(force) {
		fetch(`${Client.apiUrl}/items/refresh-file-metadata`, { method: 'POST' });
	}

	static runSpectagger() {
		fetch(`${Client.apiUrl}/spectagger/run`, { method: 'POST' });
	}

	static runItemOptimizer() {
		fetch(`${Client.apiUrl}/itemsoptimizer/run`, { method: 'POST' });
	}

	static runDirectoryScan() {
		fetch(`${Client.apiUrl}/directories/scan`, { method: 'POST' });
	}

	static getExportMetadataUrl() {
		return `${Client.apiUrl}/export-metadata.json`;
	}

	static buildFileUrl(storagePath, nonce) {
		if (!storagePath) {
			return '';
		}

		if (nonce) {
			return `${Client.apiUrl}/file/${encodeURIComponent(storagePath)}?nonce=${nonce}`;
		} else {
			return `${Client.apiUrl}/file/${encodeURIComponent(storagePath)}`;
		}
	}

	static buildInternalStoragePath(storagePath) {
		return `.internal-storage/${storagePath}`;
	}
}
