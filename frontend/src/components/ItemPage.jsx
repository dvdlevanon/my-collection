import styles from './ItemPage.module.css';
import { useEffect, useState } from 'react';
import { json, useParams } from 'react-router-dom';
import Player from './Player';
import ItemTags from './ItemTags';
import AddTagDialog from './AddTagDialog';
import { Typography } from '@mui/material';

function ItemPage() {
	const { itemId } = useParams();
	let [item, setItem] = useState(null);
	let [tags, setTags] = useState([]);
	let [addTagMode, setAddTagMode] = useState(false);

	useEffect(() => {
		reloadItem();
	}, []);

	useEffect(() => {
		fetch('http://localhost:8080/tags')
			.then((response) => response.json())
			.then((tags) => setTags(tags));
	}, []);

	const reloadItem = () => {
		fetch('http://localhost:8080/items/' + itemId)
			.then((response) => response.json())
			.then((item) => setItem(item));
	};

	const onAddTag = () => {
		setAddTagMode(true);
	};

	const onTagAdded = (tag) => {
		setAddTagMode(false);

		if (!item.tags) {
			item.tags = [];
		}

		item.tags.push(tag);

		fetch('http://localhost:8080/items/' + itemId, {
			method: 'POST',
			body: JSON.stringify(item),
		}).then(reloadItem);
	};

	const onTagRemoved = (tag) => {
		fetch('http://localhost:8080/items/' + itemId + '/remove-tag/' + tag.id, {
			method: 'POST',
			body: JSON.stringify(item),
		}).then(reloadItem);

		reloadItem();
	};

	return (
		<div className={styles.all}>
			<Typography variant="h5">{item ? item.title : ''}</Typography>
			<div className={styles.top}>
				{item ? <Player item={item} /> : ''}
				{item ? <ItemTags item={item} onAddTag={onAddTag} onTagRemoved={onTagRemoved} /> : ''}
			</div>
			<div className={styles.related_items}>related-items</div>
			{item ? (
				<AddTagDialog
					open={addTagMode}
					item={item}
					tags={tags}
					onTagAdded={onTagAdded}
					onClose={(e) => setAddTagMode(false)}
				/>
			) : (
				''
			)}
		</div>
	);
}

export default ItemPage;
