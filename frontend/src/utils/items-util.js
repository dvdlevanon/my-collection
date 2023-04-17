import Client from './client';

export default class ItemsUtil {
	static getCover = (item, coverNumber) => {
		if (item.main_cover_url) {
			return Client.buildFileUrl(item.main_cover_url);
		} else if (item.covers && item.covers.length > 0 && item.covers[coverNumber]) {
			return Client.buildFileUrl(item.covers[coverNumber].url);
		} else {
			return Client.buildFileUrl(Client.buildInternalStoragePath('covers/none/1.png'));
		}
	};
}
