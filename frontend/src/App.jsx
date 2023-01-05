import { useEffect, useState } from 'react';
import './App.css';
import Gallery from './components/Gallery';
import SuperTags from './components/SuperTags';
import Tags from './components/Tags';
import './styles/Tag.css';
import './styles/Gallery.css';
import ActiveTags from './components/ActiveTags';
import SidePanel from './components/SidePanel';

function App() {
	let [tags, setTags] = useState([]);
	let [items, setItems] = useState([]);
	let [selectedSuperTag, setSelectedSuperTag] = useState(null);
	let [activeTags, setActiveTags] = useState([]);
	let [selectedTags, setSelectedTags] = useState([]);
	// let [selectedItems, setSelectedItems] = useState([]);

	useEffect(() => {
		fetch('http://localhost:8080/items')
			.then((response) => response.json())
			.then((items) => setItems(items));
		fetch('http://localhost:8080/tags')
			.then((response) => response.json())
			.then((tags) => setTags(tags));
	}, []);

	// const getSeletedItems = (selectedTags) => {
	// 	let result = items.filter((item) => {
	// 		let tagsWithItem = selectedTags.filter((tag) => {
	// 			if (!tag.items) {
	// 				return false;
	// 			}

	// 			return (
	// 				tag.items.filter((cur) => {
	// 					return cur.id == item.id;
	// 				}).length > 0
	// 			);
	// 		});

	// 		return tagsWithItem.length === selectedTags.length;
	// 	});

	// 	return result;
	// };

	// const onToggleTag = (tag) => {
	// 	if (tag.selected) {
	// 		tag.selected = false;
	// 		let selected = selectedTags.filter((cur) => cur.id !== tag.id);
	// 		setSelectedTags(selected);
	// 		// setSelectedItems(() => getSeletedItems(selected));
	// 	} else {
	// 		tag.selected = true;
	// 		let selected = [...selectedTags, tag];
	// 		setSelectedTags(selected);
	// 		// setSelectedItems(() => getSeletedItems(selected));
	// 	}
	// };

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
	}

	const onTagDeactivated = (tag) => {
		updateTag(tag, (tag) => { 
			tag.active = false;
			tag.selected = false; 
		});

		setActiveTags(activeTags.filter((cur) => cur.id != tag.id ))
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
               		onTagSelected={onTagSelected} onTagDeselected={onTagDeselected} />
				<Gallery items={items} />
			</div>
		</div>
	);
}

export default App;
