import { Collapse, Stack } from '@mui/material';
import { useState } from 'react';
import { useQuery } from 'react-query';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import TagsUtil from '../../utils/tags-util';
import Categories from './Categories';
import Tags from './Tags';

function TagPicker({
	origin,
	onTagSelected,
	onDropDownToggled,
	initialTagSize,
	initialSelectedCategoryId,
	showSpecialCategories,
	tagLinkBuilder,
	setHideTopBar,
}) {
	const tagsQuery = useQuery(ReactQueryUtil.TAGS_KEY, Client.getTags);
	const titsQuery = useQuery(ReactQueryUtil.TAG_IMAGE_TYPES_KEY, Client.getTagImageTypes);

	let [selectedCategoryId, setSelectedCategoryId] = useState(initialSelectedCategoryId);
	let [hideCategories, setHideCategories] = useState(false);

	const getChildrenTags = (selectedId) => {
		let category = tagsQuery.data.find((cur) => {
			return cur.id == selectedId;
		});

		if (!category.children) {
			return [];
		}

		let children = category.children.map((tag) => {
			return tagsQuery.data.filter((cur) => {
				return cur.id == tag.id;
			})[0];
		});

		return children;
	};

	const onCategoryClicked = (category) => {
		if (selectedCategoryId == category.id) {
			setSelectedCategoryId(0);
			onDropDownToggled(false);
		} else {
			setSelectedCategoryId(category.id);
			onDropDownToggled(true);
		}
	};

	const tagSelectedHandler = (tag) => {
		setSelectedCategoryId(0);
		onDropDownToggled(false);
		setHideTopBar(false);
		setHideCategories(false);
		onTagSelected(tag);
	};

	return (
		<Stack className="tags_picker" height={selectedCategoryId > 0 ? '100%' : 'auto'}>
			{tagsQuery.isSuccess && (
				<Collapse in={!hideCategories}>
					<Stack>
						<Categories
							categories={TagsUtil.getCategories(tagsQuery.data).filter(
								(cur) => showSpecialCategories || TagsUtil.allowToAddToCategory(cur.id)
							)}
							onCategoryClicked={onCategoryClicked}
							selectedCategoryId={selectedCategoryId}
						/>
					</Stack>
				</Collapse>
			)}
			{tagsQuery.isSuccess && titsQuery.isSuccess && selectedCategoryId > 0 && (
				<Tags
					origin={origin}
					tags={getChildrenTags(selectedCategoryId)}
					tits={titsQuery.data}
					parent={tagsQuery.data.find((cur) => cur.id == selectedCategoryId)}
					initialTagSize={initialTagSize}
					tagLinkBuilder={tagLinkBuilder}
					onTagClicked={tagSelectedHandler}
					setHideCategories={(value) => {
						setHideTopBar(value);
						setHideCategories(value);
					}}
				/>
			)}
		</Stack>
	);
}

export default TagPicker;
