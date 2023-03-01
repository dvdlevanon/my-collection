import { Stack } from '@mui/material';
import React, { useState } from 'react';
import { useMutation, useQuery, useQueryClient } from 'react-query';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import GalleryFilters from './GalleryFilters';
import ItemsList from './ItemsList';

const itemSizes = {
	xs: { width: 150, height: 120 },
	s: { width: 250, height: 150 },
	m: { width: 350, height: 200 },
	l: { width: 450, height: 250 },
	xl: { width: 550, height: 300 },
};

function ItemsView({ previewMode }) {
	const queryClient = useQueryClient();
	const tagsQuery = useQuery(ReactQueryUtil.TAGS_KEY, Client.getTags);
	const itemsQuery = useQuery(ReactQueryUtil.ITEMS_KEY, Client.getItems);
	const [conditionType, setConditionType] = useState('||');
	const [itemsSize, setItemsSize] = useState(itemSizes.m);
	const [itemsSizeName, setItemsSizeName] = useState('m');
	const saveTag = useMutation(Client.saveTag);

	const changeTagState = (tag, updater) => {
		saveTag.mutate(updater(tag), {
			onSuccess: () => {
				ReactQueryUtil.updateTags(queryClient, tag.id, (currentTag) => {
					return updater(currentTag);
				});
			},
		});
	};

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

	const onTagDeactivated = (tag) => {
		changeTagState(tag, (currentTag) => {
			return {
				...currentTag,
				active: false,
				selected: false,
			};
		});
	};

	const onTagClick = (tag) => {
		changeTagState(tag, (currentTag) => {
			return {
				...currentTag,
				selected: !tag.selected,
			};
		});
	};

	const onChangeCondition = (conditionType) => {
		setConditionType(conditionType);
	};

	return (
		<Stack padding="10px">
			{tagsQuery.isSuccess && (
				<GalleryFilters
					conditionType={conditionType}
					activeTags={getActiveTags()}
					onTagClick={onTagClick}
					onTagDelete={onTagDeactivated}
					onChangeCondition={onChangeCondition}
					itemsSizeName={itemsSizeName}
					setItemsSizeName={(sizeName) => {
						setItemsSizeName(sizeName);
						setItemsSize(itemSizes[sizeName]);
					}}
				/>
			)}
			{tagsQuery.isSuccess && itemsQuery.isSuccess && (
				<ItemsList itemsSize={itemsSize} items={getSeletedItems(getSelectedTags())} previewMode={previewMode} />
			)}
		</Stack>
	);
}

export default ItemsView;
