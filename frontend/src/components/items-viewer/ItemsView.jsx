import { Stack } from '@mui/material';
import React, { useState } from 'react';
import { useMutation, useQueryClient } from 'react-query';
import AspectRatioUtil from '../../utils/aspect-ratio-util';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import GalleryFilters from './GalleryFilters';
import ItemsList from './ItemsList';
import ItemsViewControls from './ItemsViewControls';

function ItemsView({ previewMode, tagsQuery, itemsQuery }) {
	const queryClient = useQueryClient();
	const [conditionType, setConditionType] = useState('||');
	const [aspectRatio, setAspectRatio] = useState(AspectRatioUtil.asepctRatio16_9);
	const [itemsSize, setItemsSize] = useState({ width: 350, height: AspectRatioUtil.calcHeight(350, aspectRatio) });
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

	const onZoomChanged = (offset, aspectRatio) => {
		let newWidth = itemsSize.width + offset;
		setItemsSize({ width: newWidth, height: AspectRatioUtil.calcHeight(newWidth, aspectRatio) });
	};

	return (
		<Stack padding="10px">
			<Stack flexDirection="row" gap="10px">
				<ItemsViewControls
					itemsSize={itemsSize}
					onZoomChanged={(offest) => onZoomChanged(offest, aspectRatio)}
					aspectRatio={aspectRatio}
					onAspectRatioChanged={(newAspectRatio) => {
						setAspectRatio(newAspectRatio);
						onZoomChanged(0, newAspectRatio);
					}}
				/>
				{tagsQuery.isSuccess && (
					<GalleryFilters
						conditionType={conditionType}
						activeTags={getActiveTags()}
						onTagClick={onTagClick}
						onTagDelete={onTagDeactivated}
						onChangeCondition={onChangeCondition}
					/>
				)}
			</Stack>
			{tagsQuery.isSuccess && itemsQuery.isSuccess && (
				<ItemsList itemsSize={itemsSize} items={getSeletedItems(getSelectedTags())} previewMode={previewMode} />
			)}
		</Stack>
	);
}

export default ItemsView;
