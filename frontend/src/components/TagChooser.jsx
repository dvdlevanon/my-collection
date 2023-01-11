import { Dialog } from "@mui/material"
import { useEffect, useState } from "react";
import SuperTags from "./SuperTags";
import styles from "./TagChooser.module.css"
import Tags from "./Tags";

function TagChooser({item, onTagAdded}) {
    let [tags, setTags] = useState([]);

	useEffect(() => {
		if (tags.length != 0) {
			return;
		}

		fetch('http://localhost:8080/tags')
			.then((response) => response.json())
			.then((tags) => setTags(tags));
	}, []);

	const getSelectedSuperTag = () => {
		let selectedSupertTags = tags.filter((tag) => {
			return tag.selected && !tag.parentId;
		})

		if (selectedSupertTags.length > 0) {
			return selectedSupertTags[0];
		}

		return null;
	}

	const getTags = (superTag) => {
		if (!superTag.children) {
			return [];
		}

		let children = superTag.children.map((tag) => {
			return tags.filter((cur) => {
				return cur.id == tag.id;
			})[0];
		});

		return children
	};
	const onSuperTagSelected = (superTag) => {
		updateTag(superTag, (superTag) => {
			superTag.selected = true;
			return superTag;
		})
	};

	const updateTag = (tag, updater) => {
		setTags((tags) => { 
			return tags.map((cur) => {
				if (tag.id == cur.id) {
					return updater({...cur})
				}

				return cur
			})
		});
	}

	const onSuperTagDeselected = (superTag) => {
		updateTag(superTag, (superTag) => {
			superTag.selected = false;
			return superTag;
		})
	};

    return (
        <div className={styles.tag_chooser}>
            <SuperTags
				superTags={tags.filter((tag) => {
					return !tag.parentId;
				})}
				onSuperTagSelected={onSuperTagSelected}
				onSuperTagDeselected={onSuperTagDeselected}
			/>
            <div className={styles.tags_holder}>
                {getSelectedSuperTag() ? <Tags tags={getTags(getSelectedSuperTag())} onTagActivated={onTagAdded} /> : ''}
            </div>
        </div>
    )
}

export default TagChooser
