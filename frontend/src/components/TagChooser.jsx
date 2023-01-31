import { Box } from '@mui/material';
import { useState } from 'react';
import SuperTags from './SuperTags';
import Tags from './Tags';

function TagChooser({ tags, size, onTagSelected, onDropDownToggled, initialSelectedSuperTag }) {
	let [selectedSuperTag, setSelectedSuperTag] = useState(initialSelectedSuperTag);

	const getTags = (superTag) => {
		if (!superTag.children) {
			return [];
		}

		let children = superTag.children.map((tag) => {
			return tags.filter((cur) => {
				return cur.id == tag.id;
			})[0];
		});

		return children;
	};

	const onSuperTagClicked = (superTag) => {
		if (selectedSuperTag == superTag) {
			setSelectedSuperTag(null);
			onDropDownToggled(false);
		} else {
			setSelectedSuperTag(superTag);
			onDropDownToggled(true);
		}
	};

	const tagSelectedHandler = (tag) => {
		setSelectedSuperTag(null);
		onDropDownToggled(false);
		onTagSelected(tag);
	};

	return (
		<Box>
			<SuperTags
				superTags={tags.filter((tag) => {
					return !tag.parentId;
				})}
				onSuperTagClicked={onSuperTagClicked}
			/>
			<Box sx={{ position: 'relative' }}>
				{selectedSuperTag && (
					<Tags
						tags={getTags(selectedSuperTag)}
						parentId={selectedSuperTag.id}
						size={size}
						onTagSelected={tagSelectedHandler}
					/>
				)}
			</Box>
		</Box>
	);
}

export default TagChooser;
