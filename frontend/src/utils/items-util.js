import Client from './client';

export default class ItemsUtil {
	static PREVIEW_FROM_START_POSITION = 'start-position'; //items.go

	static getCover = (item, coverNumber) => {
		if (item.main_cover_url) {
			return Client.buildFileUrl(item.main_cover_url);
		} else if (item.covers && item.covers.length > 0 && item.covers[coverNumber]) {
			return Client.buildFileUrl(item.covers[coverNumber].url);
		} else {
			return Client.buildFileUrl(Client.buildInternalStoragePath('covers/none/1.png'));
		}
	};

	static hasPreview = (item) => {
		return item.preview_mode === PREVIEW_FROM_START_POSITION || item.preview_url;
	};

	static getPreview = (item) => {
		if (item.preview_mode === PREVIEW_FROM_START_POSITION) {
			return item.url;
		} else {
			return item.preview_url;
		}
	};
}
