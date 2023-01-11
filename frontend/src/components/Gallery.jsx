import styles from './Gallery.module.css';
import { useEffect, useState } from 'react';
import { atom, useRecoilState } from 'recoil'
import ItemsList from "./Items";
import SuperTags from './SuperTags';
import Tags from './Tags';
import SidePanel from './SidePanel';

const tagsAtom = atom({
	key: "tags",
	default: []
});

const itemsAtom = atom({
	key: "items",
	default: []
});

const conditionTypeAtom = atom({
	key: "conditionType",
	default: "||"
});

function Gallery() {
    let [tags, setTags] = useRecoilState(tagsAtom);
	let [items, setItems] = useRecoilState(itemsAtom);
	let [conditionType, setConditionType] = useRecoilState(conditionTypeAtom);

	useEffect(() => {
		if (tags.length != 0) {
			return;
		}

		fetch('http://localhost:8080/tags')
			.then((response) => response.json())
			.then((tags) => setTags(tags));
	}, []);

	useEffect(() => {
		if (items.length != 0) {
			return;
		}

		fetch('http://localhost:8080/items')
			.then((response) => response.json())
			.then((items) => setItems(items));
	}, []);

	const getSelectedTags = () => {
		return tags.filter((tag) => {
			return tag.selected;
		})
	}

	const getActiveTags = () => {
		return tags.filter((tag) => {
			return tag.active;
		})
	}

	const getSelectedSuperTag = () => {
		let selectedSupertTags = tags.filter((tag) => {
			return tag.selected && !tag.parentId;
		})

		if (selectedSupertTags.length > 0) {
			return selectedSupertTags[0];
		}

		return null;
	}

	const getSeletedItems = (selectedTags) => {
		if (selectedTags.length === 0) {
			return []
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

			if (conditionType == "&&") {
				return tagsWithItem.length === selectedTags.length;
			} else {
				return tagsWithItem.length > 0;
			}
		});

		return result;
	};

	const onSuperTagSelected = (superTag) => {
		updateTag(superTag, (superTag) => {
			superTag.selected = true;
			return superTag;
		})
	};

	const onSuperTagDeselected = (superTag) => {
		updateTag(superTag, (superTag) => {
			superTag.selected = false;
			return superTag;
		})
	};

	const onTagActivated = (tag) => {
		onSuperTagDeselected(getSelectedSuperTag())

		if (tag.active) {
			return
		}

		updateTag(tag, (tag) => { 
			tag.active = true;
			tag.selected = true;
			return tag;
		});
	}

	const onTagDeactivated = (tag) => {
		updateTag(tag, (tag) => { 
			tag.active = false;
			tag.selected = false; 
			return tag;
		});
	}
	
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

	const onTagSelected = (tag) => {
		updateTag(tag, (tag) => { 
			tag.selected = true 
			return tag;
		})
	}

	const onTagDeselected = (tag) => {
		updateTag(tag, (tag) => { 
			tag.selected = false 
			return tag;
		})
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

	const onChangeCondition = (conditionType) => {
		setConditionType(conditionType)
	};

    return (
        <div>
			<SuperTags
				superTags={tags.filter((tag) => {
					return !tag.parentId;
				})}
				onSuperTagSelected={onSuperTagSelected}
				onSuperTagDeselected={onSuperTagDeselected}
			/>
			<div className={styles.gallery_center}>
				{getSelectedSuperTag() ? <Tags tags={getTags(getSelectedSuperTag())} onTagActivated={onTagActivated} /> : ''}
				<SidePanel activeTags={getActiveTags()} onTagDeactivated={onTagDeactivated} 
					   onTagSelected={onTagSelected} onTagDeselected={onTagDeselected} onChangeCondition={onChangeCondition} />
				<ItemsList items={getSeletedItems(getSelectedTags())} />
			</div>
		</div>
    )
}

export default Gallery