export default class TasksUtil {
	static isProcessing = (task) => {
		if (TasksUtil.isDone(task)) {
			return false;
		}

		return task.processingStart != undefined && task.processingStart > 0;
	};

	static isDone = (task) => {
		return task.processingEnd != undefined;
	};

	static isPending = (task) => {
		return !TasksUtil.isDone(task) && !TasksUtil.isProcessing(task);
	};

	static isStaleProcessing = (task) => {
		if (task.processingEnd != undefined) {
			return false;
		}

		if (!TasksUtil.isProcessing(task)) {
			return false;
		}

		let millisSinceStart = Date.now() - task.processingStart;
		return millisSinceStart > 1000 * 60 * 60;
	};
}
