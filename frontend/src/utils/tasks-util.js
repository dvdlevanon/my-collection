export default class TasksUtil {
	static isProcessing = (task) => {
		return task.processingStart != undefined && task.processingStart > 0;
	};

	static isStaleProcessing = (task) => {
		if (!TasksUtil.isProcessing(task)) {
			return false;
		}

		let millisSinceStart = Date.now() - task.processingStart;
		return millisSinceStart > 1000 * 60 * 60;
	};
}
