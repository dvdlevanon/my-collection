import { useEffect, useState } from 'react';
import React from 'react'
import '../styles/Tag.css';
import '../styles/Items.css';
import '../styles/Pages.css';
import ItemsList from "../components/Items";
import SuperTags from '../components/SuperTags';
import Tags from '../components/Tags';
import SidePanel from '../components/SidePanel';
import { useRecoilState, atom } from 'recoil';

const tagsAtom = atom({
	key: "tags",
	default: []
});

const itemsAtom = atom({
	key: "items",
	default: []
});

const selectedSuperTagAtom = atom({
	key: "selectedSuperTag",
	default: null
});

const activeTagsAtom = atom({
	key: "activeTags",
	default: []
});

const selectedTagsAtom = atom({
	key: "selectedTags",
	default: []
});

const conditionTypeAtom = atom({
	key: "conditionType",
	default: "||"
});

function Gallery() {
    let [tags, setTags] = useRecoilState(tagsAtom);
	let [items, setItems] = useRecoilState(itemsAtom);
	let [selectedSuperTag, setSelectedSuperTag] = useRecoilState(selectedSuperTagAtom);
	let [activeTags, setActiveTags] = useRecoilState(activeTagsAtom);
	let [selectedTags, setSelectedTags] = useRecoilState(selectedTagsAtom);
	let [conditionType, setConditionType] = useRecoilState(conditionTypeAtom);

	useEffect(() => {
		if (items.length == 0) {
			fetch('http://localhost:8080/items')
				.then((response) => response.json())
				.then((items) => setItems(items));
		}
		if (tags.length == 0) {
			fetch('http://localhost:8080/tags')
				.then((response) => response.json())
				.then((tags) => setTags(tags));
		}
	}, []);

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
		superTag = updateTag(superTag, (tag) => { 
			tag.selected = true;
			return tag;
		});

		setSelectedSuperTag(superTag);
	};

	const onSuperTagDeselected = (superTag) => {
		superTag = updateTag(superTag, (tag) => { 
			tag.selected = false;
			return tag;
		});

		setSelectedSuperTag(null)
	};

	const onTagActivated = (tag) => {
		onSuperTagDeselected(selectedSuperTag)

		if (tag.active) {
			return
		}

		tag = updateTag(tag, (tag) => { 
			tag.active = true;
			tag.selected = true;
			return tag;
		});

		setActiveTags([...activeTags, tag])
		setSelectedTags([...selectedTags, tag])
	}

	const onTagDeactivated = (tag) => {
		tag = updateTag(tag, (tag) => { 
			tag.active = false;
			tag.selected = false;
			return tag;
		});

		setActiveTags(activeTags.filter((cur) => cur.id != tag.id ))
		setSelectedTags(selectedTags.filter((cur) => cur.id != tag.id ))
	}
	
	const updateTag = (tag, updater) => {
		let newTags = tags.map((cur) => {
			if (tag.id == cur.id) {
				return updater({...cur})
			}

			return cur
		})

		setTags(newTags);
		return newTags.filter((cur) => cur.id == tag.id)[0]
	}

	const onTagSelected = (tag) => {
		tag = updateTag(tag, (tag) => { 
			tag.selected = true;
			return tag;
		})
		setSelectedTags([...selectedTags, tag])
	}

	const onTagDeselected = (tag) => {
		tag = updateTag(tag, (tag) => {
			tag.selected = false;
			return tag;
		})
		setSelectedTags(selectedTags.filter((cur) => cur.id != tag.id ))
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
        <div className="center">
			<SuperTags
				superTags={tags.filter((tag) => {
					return !tag.parentId;
				})}
				onSuperTagSelected={onSuperTagSelected}
				onSuperTagDeselected={onSuperTagDeselected}
			/>
			<div className="center-content">
				{selectedSuperTag ? <Tags tags={getTags(selectedSuperTag)} onTagActivated={onTagActivated} /> : ''}
				<SidePanel activeTags={activeTags} onTagDeactivated={onTagDeactivated} 
					   onTagSelected={onTagSelected} onTagDeselected={onTagDeselected} onChangeCondition={onChangeCondition} />
				<ItemsList items={getSeletedItems(selectedTags)} />
			</div>
		</div>	
    )
}

export default Gallery