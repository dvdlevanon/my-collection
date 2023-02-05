import { Stack } from '@mui/material';
import { useState } from 'react';
import { useQuery } from 'react-query';
import Client from '../network/client';
import ReactQueryUtil from '../utils/react-query-util';
import SuperTags from './SuperTags';
import Tags from './Tags';

function TagChooser({ size, onTagSelected, onDropDownToggled, initialSelectedSuperTagId }) {
	const tagsQuery = useQuery(ReactQueryUtil.TAGS_KEY, Client.getTags);

	let [selectedSuperTagId, setSelectedSuperTagId] = useState(initialSelectedSuperTagId);

	const getChildrenTags = (selectedId) => {
		let superTag = tagsQuery.data.find((cur) => {
			return cur.id == selectedId;
		});

		if (!superTag.children) {
			return [];
		}

		let children = superTag.children.map((tag) => {
			return tagsQuery.data.filter((cur) => {
				return cur.id == tag.id;
			})[0];
		});

		return children;
	};

	const onSuperTagClicked = (superTag) => {
		if (selectedSuperTagId == superTag.id) {
			setSelectedSuperTagId(0);
			onDropDownToggled(false);
		} else {
			setSelectedSuperTagId(superTag.id);
			onDropDownToggled(true);
		}
	};

	const tagSelectedHandler = (tag) => {
		setSelectedSuperTagId(0);
		onDropDownToggled(false);
		onTagSelected(tag);
	};

	return (
		<Stack flexGrow={1} height="100%">
			{tagsQuery.isSuccess && (
				<SuperTags
					superTags={tagsQuery.data.filter((tag) => {
						return !tag.parentId;
					})}
					onSuperTagClicked={onSuperTagClicked}
				/>
			)}
			{tagsQuery.isSuccess && selectedSuperTagId > 0 && (
				<Tags
					tags={getChildrenTags(selectedSuperTagId)}
					parentId={selectedSuperTagId}
					size={size}
					onTagSelected={tagSelectedHandler}
				/>
			)}
		</Stack>
	);
}

export default TagChooser;
