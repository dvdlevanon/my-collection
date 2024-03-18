export default class TimeUtil {
	static msToTime = (millis) => {
		let seconds = (millis / 1000).toFixed(1);
		let minutes = (millis / (1000 * 60)).toFixed(1);
		let hours = (millis / (1000 * 60 * 60)).toFixed(1);
		let days = (millis / (1000 * 60 * 60 * 24)).toFixed(1);
		if (seconds < 1) return ' Less than a second';
		if (seconds < 60) return Math.floor(seconds) + ' Seconds';
		else if (minutes < 60) return Math.floor(minutes) + ' Minutes';
		else if (hours < 24) return Math.floor(hours) + ' Hours';
		else return Math.floor(days) + ' Days';
	};

	static formatDuration = (duration_seconds) => {
		if (!duration_seconds) {
			return '00:00';
		}

		if (duration_seconds < 60 * 60) {
			return new Date(duration_seconds * 1000).toISOString().slice(14, 19);
		} else {
			return new Date(duration_seconds * 1000).toISOString().slice(11, 19);
		}
	};

	static formatEpochToDate(epoch) {
		var date = new Date(epoch);
		var formattedDate =
			date.getFullYear() +
			'/' +
			(date.getMonth() + 1).toString().padStart(2, '0') +
			'/' +
			date.getDate().toString().padStart(2, '0') +
			' ' +
			date.getHours().toString().padStart(2, '0') +
			':' +
			date.getMinutes().toString().padStart(2, '0') +
			':' +
			date.getSeconds().toString().padStart(2, '0');

		return formattedDate;
	}

	static formatSeconds(seconds) {
		return (
			Math.floor(seconds / 60)
				.toString()
				.padStart(2, '0') +
			':' +
			Math.floor(seconds % 60)
				.toString()
				.padStart(2, '0')
		);
	}
}
