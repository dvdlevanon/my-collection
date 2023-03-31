import { Divider, Stack } from '@mui/material';
import { useEffect, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from 'react-query';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import ItemsView from '../items-viewer/ItemsView';
import TagPicker from '../tag-picker/TagPicker';

function Gallery({ previewMode }) {
	const queryClient = useQueryClient();
	const tagsQuery = useQuery(ReactQueryUtil.TAGS_KEY, Client.getTags);
	const itemsQuery = useQuery(ReactQueryUtil.ITEMS_KEY, Client.getItems);
	const saveTag = useMutation(Client.saveTag);
	let [tagsDropDownOpened, setTagsDropDownOpened] = useState(false);

	useEffect(() => {
		document.title = 'My Collection';
	}, []);

	const changeTagState = (tag, updater) => {
		saveTag.mutate(updater(tag), {
			onSuccess: () => {
				ReactQueryUtil.updateTags(queryClient, tag.id, (currentTag) => {
					return updater(currentTag);
				});
			},
		});
	};

	const onTagActivated = (tag) => {
		window.scrollTo(0, 0);
		changeTagState(tag, (currentTag) => {
			return {
				...currentTag,
				active: true,
				selected: true,
			};
		});
	};

	return (
		<Stack flexGrow={1} padding="10px">
			{tagsQuery.isSuccess && (
				<TagPicker
					size="big"
					showSpecialCategories={true}
					onTagSelected={onTagActivated}
					onDropDownToggled={(state) => setTagsDropDownOpened(state)}
					initialSelectedCategoryId={0}
				/>
			)}
			<Divider />
			<Stack padding="10px">
				{!tagsDropDownOpened && (
					<ItemsView tagsQuery={tagsQuery} itemsQuery={itemsQuery} previewMode={previewMode} />
				)}
			</Stack>
		</Stack>
	);
}

export default Gallery;
