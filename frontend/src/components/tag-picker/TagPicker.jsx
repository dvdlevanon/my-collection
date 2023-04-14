import { Stack } from '@mui/material';
import { useState } from 'react';
import { useQuery } from 'react-query';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import TagsUtil from '../../utils/tags-util';
import Categories from './Categories';
import Tags from './Tags';

function TagPicker({
	size,
	onTagSelected,
	onDropDownToggled,
	initialSelectedCategoryId,
	showSpecialCategories,
	tagLinkBuilder,
}) {
	const tagsQuery = useQuery(ReactQueryUtil.TAGS_KEY, Client.getTags);

	let [selectedCategoryId, setSelectedCategoryId] = useState(initialSelectedCategoryId);

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
		onTagSelected(tag);
	};

	return (
		<Stack height={selectedCategoryId > 0 ? '100%' : 'auto'}>
			{tagsQuery.isSuccess && (
				<Categories
					categories={TagsUtil.getCategories(tagsQuery.data).filter(
						(cur) => showSpecialCategories || !TagsUtil.isSpecialCategory(cur.id)
					)}
					onCategoryClicked={onCategoryClicked}
					selectedCategoryId={selectedCategoryId}
				/>
			)}
			{tagsQuery.isSuccess && selectedCategoryId > 0 && (
				<Tags
					tags={getChildrenTags(selectedCategoryId)}
					parentId={selectedCategoryId}
					size={size}
					tagLinkBuilder={tagLinkBuilder}
					onTagSelected={tagSelectedHandler}
				/>
			)}
		</Stack>
	);
}

export default TagPicker;
