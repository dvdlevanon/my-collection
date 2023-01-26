import { Divider } from '@mui/material';
import { useEffect, useState } from 'react';
import Client from '../network/client';
import GalleryFilters from './GalleryFilters';
import ItemsList from './Items';
import TagChooser from './TagChooser';

function Gallery({ previewMode }) {
	let [tags, setTags] = useState([]);
	let [items, setItems] = useState([]);
	let [conditionType, setConditionType] = useState('||');
	let [tagsDropDownOpened, setTagsDropDownOpened] = useState(false);

	useEffect(() => Client.getTags((tags) => setTags(tags)), []);
	useEffect(() => Client.getItems((items) => setItems(items)), []);

	const getSelectedTags = () => {
		return tags.filter((tag) => {
			return tag.selected;
		});
	};

	const getActiveTags = () => {
		return tags.filter((tag) => {
			return tag.active;
		});
	};

	const getSeletedItems = (selectedTags) => {
		if (selectedTags.length === 0) {
			return [];
		}

		let result = items.filter((item) => {
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

	const onTagActivated = (tag) => {
		tag.active = true;
		tag.selected = true;
		Client.saveTag(tag, () => Client.getTags((tags) => setTags(tags)));
	};

	const onTagDeactivated = (tag) => {
		tag.active = false;
		tag.selected = false;
		Client.saveTag(tag, () => Client.getTags((tags) => setTags(tags)));
	};

	const onTagSelected = (tag) => {
		tag.selected = true;
		Client.saveTag(tag, () => Client.getTags((tags) => setTags(tags)));
	};

	const onTagDeselected = (tag) => {
		tag.selected = false;
		Client.saveTag(tag, () => Client.getTags((tags) => setTags(tags)));
	};

	const onChangeCondition = (conditionType) => {
		setConditionType(conditionType);
	};

	return (
		<div>
			<TagChooser
				tags={tags}
				markActive={true}
				onTagSelected={onTagActivated}
				onDropDownToggled={(state) => setTagsDropDownOpened(state)}
			/>
			<Divider sx={{ borderBottomWidth: 2 }} />
			{!tagsDropDownOpened && (
				<GalleryFilters
					activeTags={getActiveTags()}
					onTagDeactivated={onTagDeactivated}
					onTagSelected={onTagSelected}
					onTagDeselected={onTagDeselected}
					onChangeCondition={onChangeCondition}
				/>
			)}
			{!tagsDropDownOpened && <ItemsList items={getSeletedItems(getSelectedTags())} previewMode={previewMode} />}
		</div>
	);
}

export default Gallery;
