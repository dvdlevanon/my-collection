import { Stack } from '@mui/material';
import React, { useState } from 'react';
import AspectRatioUtil from '../../utils/aspect-ratio-util';
import GalleryFilters from './GalleryFilters';
import ItemsList from './ItemsList';
import ItemsViewControls from './ItemsViewControls';

function ItemsView({ previewMode, tagsQuery, itemsQuery, galleryUrlParams }) {
	const [conditionType, setConditionType] = useState('||');
	const [searchTerm, setSearchTerm] = useState('');
	const [aspectRatio, setAspectRatio] = useState(AspectRatioUtil.asepctRatio16_9);
	const [itemsSize, setItemsSize] = useState({ width: 350, height: AspectRatioUtil.calcHeight(350, aspectRatio) });

	const getSelectedTags = () => {
		let selectedTagsIds = galleryUrlParams.getSelectedTags();

		return tagsQuery.data.filter((tag) => {
			return selectedTagsIds.some((id) => tag.id == id);
		});
	};

	const getActiveTags = () => {
		let activeTagsIds = galleryUrlParams.getActiveTags();

		return tagsQuery.data.filter((tag) => {
			return activeTagsIds.some((id) => tag.id == id);
		});
	};

	const isMetSearchTerm = (item) => {
		return item.title && item.title.toLowerCase().includes(searchTerm.toLowerCase());
	};

	const getFilteredItems = (selectedTags, searchTerm) => {
		if (selectedTags.length === 0) {
			if (!searchTerm || searchTerm.length < 3) {
				return [];
			}

			let filtered = itemsQuery.data.filter((item) => isMetSearchTerm(item));
			if (filtered.length > 500) {
				return filtered.slice(500);
			}
			return filtered;
		}

		let result = itemsQuery.data.filter((item) => {
			if (searchTerm && !isMetSearchTerm(item)) {
				return false;
			}

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
		galleryUrlParams.deactivateTag(tag.id);
	};

	const onTagClick = (tag) => {
		galleryUrlParams.toggleTagSelection(tag.id);
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
						selectedTags={getSelectedTags()}
						onTagClick={onTagClick}
						onTagDelete={onTagDeactivated}
						onChangeCondition={onChangeCondition}
						searchTerm={searchTerm}
						setSearchTerm={setSearchTerm}
						galleryUrlParams={galleryUrlParams}
					/>
				)}
			</Stack>
			{tagsQuery.isSuccess && itemsQuery.isSuccess && (
				<ItemsList
					itemsSize={itemsSize}
					items={getFilteredItems(getSelectedTags(), searchTerm)}
					previewMode={previewMode}
					itemLinkBuilder={(item) => {
						return '/spa/item/' + item.id + '?' + galleryUrlParams.getUrlParamsString();
					}}
				/>
			)}
		</Stack>
	);
}

export default ItemsView;
