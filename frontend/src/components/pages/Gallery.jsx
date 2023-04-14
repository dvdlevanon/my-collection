import { Divider, Stack } from '@mui/material';
import { useEffect, useState } from 'react';
import { useQuery } from 'react-query';
import { useSearchParams } from 'react-router-dom';
import Client from '../../utils/client';
import GalleryUrlParams from '../../utils/gallery-url-params';
import ReactQueryUtil from '../../utils/react-query-util';
import ItemsView from '../items-viewer/ItemsView';
import TagPicker from '../tag-picker/TagPicker';

function Gallery({ previewMode }) {
	const tagsQuery = useQuery(ReactQueryUtil.TAGS_KEY, Client.getTags);
	const itemsQuery = useQuery(ReactQueryUtil.ITEMS_KEY, Client.getItems);
	const [tagsDropDownOpened, setTagsDropDownOpened] = useState(false);
	let [searchParams, setSearchParams] = useSearchParams();
	let galleryUrlParams = new GalleryUrlParams(searchParams, setSearchParams);

	useEffect(() => {
		document.title = 'My Collection';
	}, []);

	const onTagActivated = (tag) => {
		galleryUrlParams.activateTag(tag.id);
		window.scrollTo(0, 0);
	};

	return (
		<Stack flexGrow={1} padding="10px">
			<Stack>
				{tagsQuery.isSuccess && (
					<TagPicker
						size="big"
						showSpecialCategories={true}
						onTagSelected={onTagActivated}
						onDropDownToggled={(state) => setTagsDropDownOpened(state)}
						initialSelectedCategoryId={0}
						tagLinkBuilder={(tag) => {
							return '/?' + galleryUrlParams.buildActivateTagUrl(tag.id);
						}}
					/>
				)}
			</Stack>
			<Divider />
			<Stack padding="10px">
				{!tagsDropDownOpened && (
					<ItemsView
						tagsQuery={tagsQuery}
						itemsQuery={itemsQuery}
						galleryUrlParams={galleryUrlParams}
						previewMode={previewMode}
					/>
				)}
			</Stack>
		</Stack>
	);
}

export default Gallery;
