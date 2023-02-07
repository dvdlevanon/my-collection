import { Box, Stack, Typography } from '@mui/material';
import { useState } from 'react';
import { useQuery, useQueryClient } from 'react-query';
import { useParams } from 'react-router-dom';
import Client from '../network/client';
import ReactQueryUtil from '../utils/react-query-util';
import AttachTagDialog from './AttachTagDialog';
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
		<Stack>
			{itemQuery.isSuccess && <Typography variant="h5">{itemQuery.data.title}</Typography>}
			{itemQuery.isSuccess && (
				<Stack flexDirection="row">
					<Box flexGrow={1} padding="10px" height="50%">
						<Player item={itemQuery.data} />
					</Box>
					<ItemTags item={itemQuery.data} onAddTag={onAddTag} onTagRemoved={onTagRemoved} />
				</Stack>
			)}
			{itemQuery.isSuccess && tagsQuery.isSuccess && (
				<AttachTagDialog
					open={addTagMode}
					item={itemQuery.data}
					onTagAdded={onTagAdded}
					onClose={(e) => setAddTagMode(false)}
				/>
			)}
		</Stack>
	);
}

export default ItemPage;
