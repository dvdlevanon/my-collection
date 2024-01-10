import { Box, Fade, Stack } from '@mui/material';
import React, { useEffect, useState } from 'react';
import seedrandom from 'seedrandom';
import AspectRatioUtil from '../../utils/aspect-ratio-util';
import TagsUtil from '../../utils/tags-util';
import ManageTagImageDialog from '../dialogs/ManageTagImageDialog';
import TagThumbnails from '../tag-thumbnail/TagThumbnails';
import GalleryFilters from './GalleryFilters';
import ItemSortSelector from './ItemSortSelector';
import ItemsList from './ItemsList';
import ItemsViewControls from './ItemsViewControls';

function ItemsView({ previewMode, tagsQuery, itemsQuery, galleryUrlParams }) {
	const [conditionType, setConditionType] = useState(galleryUrlParams.getConditionType());
	const [searchTerm, setSearchTerm] = useState('');
	const [aspectRatio, setAspectRatio] = useState(AspectRatioUtil.asepctRatio16_9);
	const [itemsSize, setItemsSize] = useState({ width: 350, height: AspectRatioUtil.calcHeight(350, aspectRatio) });
	const [sortBy, setSortBy] = useState('random');
	const [editThumbnailTag, setEditThumbnailTag] = useState(null);
	const [showThumbnails, setShowThumbnails] = useState(true);

	useEffect(() => {
		let lastItemsWidth = localStorage.getItem('items-width');
		if (lastItemsWidth) {
			let lastItemsWidthInt = parseInt(lastItemsWidth);
			setItemsSize({
				width: lastItemsWidthInt,
				height: AspectRatioUtil.calcHeight(lastItemsWidthInt, aspectRatio),
			});
		}

		let lastSortBy = localStorage.getItem('items-sort-by');
		if (lastSortBy) {
			setSortBy(lastSortBy);
		} else {
			setSortBy('random');
		}
	}, []);

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

	const sortItems = (items) => {
		if (sortBy == 'random') {
			let epochDay = Math.floor(Date.now() / 1000 / 60 / 60 / 24);
			let randomItems = [];
			let rand = seedrandom(epochDay);

			for (let i = 0; i < items.length; i++) {
				let randomIndex = Math.floor(rand() * items.length);
				while (randomItems[randomIndex]) {
					randomIndex = Math.floor(rand() * items.length);
				}

				randomItems[randomIndex] = items[i];
			}

			return randomItems;
		} else if (sortBy == 'title-asc') {
			return items.sort((a, b) => (a.title > b.title ? 1 : a.title < b.title ? -1 : 0));
		} else if (sortBy == 'title-desc') {
			return items.sort((a, b) => (a.title > b.title ? -1 : a.title < b.title ? 1 : 0));
		} else if (sortBy == 'duration-desc') {
			return items.sort((a, b) =>
				a.duration_seconds > b.duration_seconds ? 1 : a.duration_seconds < b.duration_seconds ? -1 : 0
			);
		} else if (sortBy == 'duration-asc' || sortBy == 'duration') {
			return items.sort((a, b) =>
				a.duration_seconds > b.duration_seconds ? -1 : a.duration_seconds < b.duration_seconds ? 1 : 0
			);
		} else {
			return items;
		}
	};

	const getAvailableThumbnailTags = (items) => {
		let tags = [];
		let selectedTags = getSelectedTags();

		for (let i = 0; i < items.length; i++) {
			for (let j = 0; j < items[i].tags.length; j++) {
				let tag = items[i].tags[j];
				if (!TagsUtil.showAsThumbnail(tag.parentId)) {
					continue;
				}

				if (selectedTags.some((cur) => tag.id == cur.id)) {
					continue;
				}

				tags[tag.id] = tag;
			}
		}

		let result = [];

		for (let i = 0; i < Object.keys(tags).length; i++) {
			result.push(tags[Object.keys(tags)[i]]);
		}

		return result.sort((a, b) => (a.parentId > b.parentId ? 1 : a.parentId < b.parentId ? 1 : 0));
	};

	const onTagDeactivated = (tag) => {
		galleryUrlParams.deactivateTag(tag.id);
	};

	const onTagClick = (tag) => {
		galleryUrlParams.toggleTagSelection(tag.id);
	};

	const onZoomChanged = (offset, aspectRatio) => {
		let newWidth = itemsSize.width + offset;
		localStorage.setItem('items-width', newWidth);
		setItemsSize({ width: newWidth, height: AspectRatioUtil.calcHeight(newWidth, aspectRatio) });
	};

	const onSortChanged = (newSortBy) => {
		if (newSortBy == 'duration' && sortBy.startsWith('duration')) {
			newSortBy = sortBy == 'duration-desc' ? 'duration-asc' : 'duration-desc';
		}

		if (sortBy == newSortBy) {
			return;
		}

		localStorage.setItem('items-sort-by', newSortBy);
		setSortBy(newSortBy);
	};

	const updateConditionType = (type) => {
		galleryUrlParams.setConditionType(type);
		setConditionType(type);
	};

	return (
		<Stack overflow="hidden" height="100%">
			<Stack flexDirection="row" gap="10px" padding="0px 0px 3px 0px">
				<ItemsViewControls
					itemsSize={itemsSize}
					onZoomChanged={(offest) => onZoomChanged(offest, aspectRatio)}
					aspectRatio={aspectRatio}
					onAspectRatioChanged={(newAspectRatio) => {
						setAspectRatio(newAspectRatio);
						onZoomChanged(0, newAspectRatio);
					}}
				/>
				<ItemSortSelector sortBy={sortBy} onSortChanged={onSortChanged} />
				{tagsQuery.isSuccess && (
					<GalleryFilters
						conditionType={conditionType}
						setConditionType={updateConditionType}
						activeTags={getActiveTags()}
						selectedTags={getSelectedTags()}
						onTagClick={onTagClick}
						onTagDelete={onTagDeactivated}
						searchTerm={searchTerm}
						setSearchTerm={setSearchTerm}
						galleryUrlParams={galleryUrlParams}
					/>
				)}
			</Stack>
			{tagsQuery.isSuccess && itemsQuery.isSuccess && (
				<Fade in={showThumbnails} unmountOnExit={true}>
					<Box>
						<TagThumbnails
							tags={getAvailableThumbnailTags(getFilteredItems(getSelectedTags(), searchTerm))}
							onEditThumbnail={(tag) => setEditThumbnailTag(tag)}
							onTagClicked={(tag) => {
								updateConditionType('&&');
								galleryUrlParams.activateTag(tag.id, true);
							}}
							carouselMode={true}
						/>
					</Box>
				</Fade>
			)}
			<Box width="100%" height="100%">
				{tagsQuery.isSuccess && itemsQuery.isSuccess && (
					<ItemsList
						itemsSize={itemsSize}
						items={sortItems(getFilteredItems(getSelectedTags(), searchTerm))}
						previewMode={previewMode}
						itemLinkBuilder={(item) => {
							return '/spa/item/' + item.id + '?' + galleryUrlParams.getUrlParamsString();
						}}
						onScroll={(e) => {
							if (!e) {
								return;
							}
							setShowThumbnails(e.scrollTop < 10);
						}}
					/>
				)}
			</Box>
			{editThumbnailTag != null && (
				<ManageTagImageDialog
					tag={editThumbnailTag}
					autoThumbnailMode={true}
					onClose={() => setEditThumbnailTag(null)}
				/>
			)}
		</Stack>
	);
}

export default ItemsView;
