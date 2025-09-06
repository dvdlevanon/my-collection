export default class ColorUtil {
	static getInvertedColor = (color) => {
		if (!color) {
			return color;
		}

		if (!color.startsWith('#')) {
			return color; // only hex color supported
		}

		const hex = color.replace('#', '');
		const r = parseInt(hex.substr(0, 2), 16);
		const g = parseInt(hex.substr(2, 2), 16);
		const b = parseInt(hex.substr(4, 2), 16);

		return `rgb(${255 - r}, ${255 - g}, ${255 - b})`;
	};
}
