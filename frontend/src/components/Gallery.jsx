import { Divider } from '@mui/material';
import { useState } from 'react';
import { useMutation, useQuery, useQueryClient } from 'react-query';
import Client from '../network/client';
import ReactQueryUtil from '../utils/react-query-util';
import GalleryFilters from './GalleryFilters';
import ItemsList from './Items';
import TagChooser from './TagChooser';

function Gallery({ previewMode }) {
	const queryClient = useQueryClient();
	const tagsQuery = useQuery(ReactQueryUtil.TAGS_KEY, Client.getTags);
	const itemsQuery = useQuery(ReactQueryUtil.ITEMS_KEY, Client.getItems);
	const saveTag = useMutation(Client.saveTag);
	let [conditionType, setConditionType] = useState('||');
	let [tagsDropDownOpened, setTagsDropDownOpened] = useState(false);

	const getSelectedTags = () => {
		return tagsQuery.data.filter((tag) => {
			return tag.selected;
		});
	};

	const getActiveTags = () => {
		return tagsQuery.data.filter((tag) => {
			return tag.active;
		});
	};

	const getSeletedItems = (selectedTags) => {
		if (selectedTags.length === 0) {
			return [];
		}

		let result = itemsQuery.data.filter((item) => {
			let tagsWithItem = selectedTags.filter((tag) => {
				if (!tag.items) {
					return false;
				}

				return (
					tag.items.filter((cur) => {
						return cur.id == item.id;
					}).length > 0
				);
			});

			if (conditionType == '&&') {
				return tagsWithItem.length === selectedTags.length;
			} else {
				return tagsWithItem.length > 0;
			}
		});

		return result;
	};

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
		changeTagState(tag, (currentTag) => {
			return {
				...currentTag,
				active: true,
				selected: true,
			};
		});
	};

	const onTagDeactivated = (tag) => {
		changeTagState(tag, (currentTag) => {
			return {
				...currentTag,
				active: false,
				selected: false,
			};
		});
	};

	const onTagSelected = (tag) => {
		changeTagState(tag, (currentTag) => {
			return {
				...currentTag,
				selected: true,
			};
		});
	};

	const onTagDeselected = (tag) => {
		changeTagState(tag, (currentTag) => {
			return {
				...currentTag,
				selected: false,
			};
		});
	};

	const onChangeCondition = (conditionType) => {
		setConditionType(conditionType);
	};

	return (
		<div>
			{tagsQuery.isSuccess && (
				<TagChooser
					tags={tagsQuery.data}
					size="big"
					onTagSelected={onTagActivated}
					onDropDownToggled={(state) => setTagsDropDownOpened(state)}
				/>
			)}
			<Divider sx={{ borderBottomWidth: 2 }} />
			{tagsQuery.isSuccess && !tagsDropDownOpened && (
				<GalleryFilters
					activeTags={getActiveTags()}
					onTagDeactivated={onTagDeactivated}
					onTagSelected={onTagSelected}
					onTagDeselected={onTagDeselected}
					onChangeCondition={onChangeCondition}
				/>
			)}
			{tagsQuery.isSuccess && itemsQuery.isSuccess && !tagsDropDownOpened && (
				<ItemsList items={getSeletedItems(getSelectedTags())} previewMode={previewMode} />
			)}
		</div>
	);
}

export default Gallery;
