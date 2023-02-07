import AddIcon from '@mui/icons-material/Add';
import { Box, Chip, Container, Stack, Typography } from '@mui/material';
import { useState } from 'react';
import { useQuery, useQueryClient } from 'react-query';
import { useParams } from 'react-router-dom';
import Client from '../network/client';
import ReactQueryUtil from '../utils/react-query-util';
import AttachTagDialog from './AttachTagDialog';
import Player from './Player';
import TagChips from './TagChips';

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

	const onTagClicked = (tag) => {
		console.log('send to gallery with selected filter');
	};

	return (
		<Container maxWidth="xl">
			{itemQuery.isSuccess && (
				<Stack flexGrow={1} flexDirection="column" padding="20px">
					<Typography variant="h5">{itemQuery.data.title}</Typography>
					<Box padding="10px">
						<Player item={itemQuery.data} />
					</Box>
					<Stack flexDirection="row" gap="10px">
						<TagChips
							flexDirection="column"
							tags={itemQuery.data.tags}
							onDelete={onTagRemoved}
							onClick={onTagClicked}
							tagHighlightedPredicate={() => {
								return true;
							}}
						></TagChips>
						<Chip
							color="secondary"
							icon={<AddIcon />}
							onClick={onAddTag}
							sx={{ '& .MuiChip-label': { padding: '5px' } }}
						/>
					</Stack>
					<AttachTagDialog
						open={addTagMode}
						item={itemQuery.data}
						onTagAdded={onTagAdded}
						onClose={(e) => setAddTagMode(false)}
					/>
				</Stack>
			)}
		</Container>
	);
}

export default ItemPage;
