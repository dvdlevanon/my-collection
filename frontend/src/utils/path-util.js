export default class PathUtil {
	static dirname = (path) => {
		const parts = path.split('/');
		return parts[parts.length - 1] || '.';
	};
}
