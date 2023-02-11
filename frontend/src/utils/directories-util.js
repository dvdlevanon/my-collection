export default class DirectoriesUtil {
	static isProcessing = (directory) => {
		return directory.processingStart != undefined && directory.processingStart > 0;
	};

	static isStaleProcessing = (directory) => {
		if (!DirectoriesUtil.isProcessing(directory)) {
			return false;
		}

		let millisSinceStart = Date.now() - directory.processingStart;
		return millisSinceStart > 1000 * 60;
	};
}
