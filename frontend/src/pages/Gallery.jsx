import { useEffect, useState } from 'react';
import React from 'react'
import '../styles/Tag.css';
import '../styles/Items.css';
import '../styles/Pages.css';
import ItemsList from "../components/Items";
import SuperTags from '../components/SuperTags';
import Tags from '../components/Tags';
import SidePanel from '../components/SidePanel';

function Gallery() {
    let [tags, setTags] = useState([]);
	let [items, setItems] = useState([]);
	let [selectedSuperTag, setSelectedSuperTag] = useState(null);
	let [activeTags, setActiveTags] = useState([]);
	let [selectedTags, setSelectedTags] = useState([]);
	let [conditionType, setConditionType] = useState("||");

	useEffect(() => {
		fetch('http://localhost:8080/items')
			.then((response) => response.json())
			.then((items) => setItems(items));
		fetch('http://localhost:8080/tags')
			.then((response) => response.json())
			.then((tags) => setTags(tags));
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
		superTag.selected = true;
		setSelectedSuperTag(superTag);
	};

	const onSuperTagDeselected = (superTag) => {
		selectedSuperTag.selected = false;
		setSelectedSuperTag(null)
	};

	const onTagActivated = (tag) => {
		onSuperTagDeselected(selectedSuperTag)

		if (tag.active) {
			return
		}

		updateTag(tag, (tag) => { 
			tag.active = true;
			tag.selected = true; 
		});

		setActiveTags([...activeTags, tag])
		setSelectedTags([...selectedTags, tag])
	}

	const onTagDeactivated = (tag) => {
		updateTag(tag, (tag) => { 
			tag.active = false;
			tag.selected = false; 
		});

		setActiveTags(activeTags.filter((cur) => cur.id != tag.id ))
		setSelectedTags(selectedTags.filter((cur) => cur.id != tag.id ))
	}
	
	const updateTag = (tag, updater) => {
		setTags(tags.map((cur) => {
			if (tag.id == cur.id) {
				updater(cur)
			}

			return cur
		}));
	}

	const onTagSelected = (tag) => {
		updateTag(tag, (tag) => { tag.selected = true })
		setSelectedTags([...selectedTags, tag])
	}

	const onTagDeselected = (tag) => {
		updateTag(tag, (tag) => { tag.selected = false })
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