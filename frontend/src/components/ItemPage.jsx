import { Typography } from '@mui/material';
import { useState } from 'react';
import { useQuery, useQueryClient } from 'react-query';
import { useParams } from 'react-router-dom';
import Client from '../network/client';
import ReactQueryUtil from '../utils/react-query-util';
import AddTagDialog from './AddTagDialog';
import styles from './ItemPage.module.css';
import ItemTags from './ItemTags';
import Player from './Player';

function ItemPage() {
	const queryClient = useQueryClient();
	const { itemId } = useParams();
	const itemQuery = useQuery(ReactQueryUtil.itemKey(itemId), () => Client.getItem(itemId));
	const tagsQuery = useQuery(ReactQueryUtil.TAGS_KEY, Client.getTags);
	let [addTagMode, setAddTagMode] = useState(false);

	const onAddTag = () => {
		setAddTagMode(true);
	};

	const onTagAdded = (tag) => {
		setAddTagMode(false);

		let tags = itemQuery.data.tags || [];
		tags.push(tag);

		Client.saveItem({ ...itemQuery.data, tags: tags }, () => {
			queryClient.refetchQueries({ queryKey: itemQuery.queryKey });
		});
	};

	const onTagRemoved = (tag) => {
		Client.removeTagFromItem(itemQuery.data.id, tag.id, () => {
			queryClient.refetchQueries({ queryKey: itemQuery.queryKey });
		});
	};

	return (
		<div className={styles.all}>
			{itemQuery.isSuccess && <Typography variant="h5">{itemQuery.data.title}</Typography>}
			{itemQuery.isSuccess && (
				<div className={styles.top}>
					<Player item={itemQuery.data} />
					<ItemTags item={itemQuery.data} onAddTag={onAddTag} onTagRemoved={onTagRemoved} />
				</div>
			)}
			{itemQuery.isSuccess && tagsQuery.isSuccess && (
				<AddTagDialog
					open={addTagMode}
					item={itemQuery.data}
					tags={tagsQuery.data}
					onTagAdded={onTagAdded}
					onClose={(e) => setAddTagMode(false)}
				/>
			)}
		</div>
	);
}

export default ItemPage;
