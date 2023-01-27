import { useState } from 'react';
import SuperTags from './SuperTags';
import styles from './TagChooser.module.css';
import Tags from './Tags';

function TagChooser({ tags, size, onTagSelected, onDropDownToggled }) {
	let [selectedSuperTag, setSelectedSuperTag] = useState(null);

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
		<div className={styles.tag_chooser}>
			<SuperTags
				superTags={tags.filter((tag) => {
					return !tag.parentId;
				})}
				onSuperTagClicked={onSuperTagClicked}
			/>
			<div className={styles.tags_holder}>
				{selectedSuperTag && (
					<Tags tags={getTags(selectedSuperTag)} size={size} onTagSelected={tagSelectedHandler} />
				)}
			</div>
		</div>
	);
}

export default TagChooser;
