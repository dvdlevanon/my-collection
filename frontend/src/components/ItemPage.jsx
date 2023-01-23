import { Typography } from '@mui/material';
import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import Client from '../network/client';
import AddTagDialog from './AddTagDialog';
import styles from './ItemPage.module.css';
import ItemTags from './ItemTags';
import Player from './Player';

function ItemPage() {
	const { itemId } = useParams();
	let [item, setItem] = useState(null);
	let [tags, setTags] = useState([]);
	let [addTagMode, setAddTagMode] = useState(false);

	useEffect(() => Client.getTags((tags) => setTags(tags)), []);
	useEffect(() => reloadItem(), []);

	const reloadItem = () => {
		Client.getItem(itemId, (item) => setItem(item));
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
		Client.saveItem(item, reloadItem);
	};

	const onTagRemoved = (tag) => {
		Client.removeTagFromItem(item.id, tag.id, reloadItem);
	};

	return (
		<div className={styles.all}>
			<Typography variant="h5">{item && item.title}</Typography>
			<div className={styles.top}>
				{item && <Player item={item} />}
				{item && <ItemTags item={item} onAddTag={onAddTag} onTagRemoved={onTagRemoved} />}
			</div>
			{item && (
				<AddTagDialog
					open={addTagMode}
					item={item}
					tags={tags}
					onTagAdded={onTagAdded}
					onClose={(e) => setAddTagMode(false)}
				/>
			)}
		</div>
	);
}

export default ItemPage;
