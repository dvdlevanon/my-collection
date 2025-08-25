import { Collapse, Stack } from '@mui/material';
import { useQuery } from '@tanstack/react-query';
import { useState } from 'react';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import TagsUtil from '../../utils/tags-util';
import DirectoriesTree from '../directories-tree/DirectoriesTree';
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
	singleCategoryMode,
	galleryUrlParams,
}) {
	const tagsQuery = useQuery({ queryKey: ReactQueryUtil.TAGS_KEY, queryFn: Client.getTags });
	const titsQuery = useQuery({ queryKey: ReactQueryUtil.TAG_IMAGE_TYPES_KEY, queryFn: Client.getTagImageTypes });

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

	const buildDirectoriesComponent = () => {
		return <DirectoriesTree />;
	};

	const buildDefaultTagsComponent = () => {
		return (
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
		);
	};

	const shouldShowTags = () => {
		return tagsQuery.isSuccess && titsQuery.isSuccess && selectedCategoryId > 0;
	};

	const buildTagsComponent = () => {
		if (TagsUtil.isDirectoriesCategory(selectedCategoryId)) {
			return buildDirectoriesComponent();
		} else {
			return buildDefaultTagsComponent();
		}
	};

	return (
		<Stack className="tags_picker" height={selectedCategoryId > 0 ? '100%' : 'auto'}>
			{tagsQuery.isSuccess && !singleCategoryMode && (
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
			{shouldShowTags() && buildTagsComponent()}
		</Stack>
	);
}

export default TagPicker;
