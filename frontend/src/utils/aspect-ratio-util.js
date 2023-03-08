export default class AspectRatioUtil {
	static asepctRatio16_9 = { x: 16, y: 9 };
	static asepctRatio4_3 = { x: 4, y: 3 };
	static asepctRatio4_2 = { x: 4, y: 2 };
	static all = [AspectRatioUtil.asepctRatio16_9, AspectRatioUtil.asepctRatio4_3, AspectRatioUtil.asepctRatio4_2];

	static calcHeight = (width, ratio) => {
		let ar = ratio.y / ratio.x;
		return width * ar;
	};

	static nearestAspectRatio = (width, height) => {
		let ar = height / width;
		let result = AspectRatioUtil.all[0];
		let min = 10000;

		for (let i = 0; i < AspectRatioUtil.all.length; i++) {
			let cur = AspectRatioUtil.all[i].y / AspectRatioUtil.all[i].x;
			let diff = Math.abs(cur - ar);

			if (diff < min) {
				min = diff;
				result = AspectRatioUtil.all[i];
			}
		}

		return result;
	};

	static toString = (ar) => {
		return ar.x + ':' + ar.y;
	};
}
