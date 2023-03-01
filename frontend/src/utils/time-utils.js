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
}
