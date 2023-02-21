import AddIcon from '@mui/icons-material/Add';
import { Box, Chip, Stack, Typography } from '@mui/material';
import { useLayoutEffect, useState } from 'react';
import { useQuery, useQueryClient } from 'react-query';
import { useParams } from 'react-router-dom';
import Client from '../../network/client';
import ReactQueryUtil from '../../utils/react-query-util';
import TagsUtil from '../../utils/tags-util';
import AttachTagDialog from '../dialogs/AttachTagDialog';
import TagChips from '../tags-chip/TagChips';
import Player from './Player';

function useWindowSize() {
	const [size, setSize] = useState([0, 0]);
	useLayoutEffect(() => {
		function updateSize() {
			setSize([window.innerWidth, window.innerHeight]);
		}
		window.addEventListener('resize', updateSize);
		updateSize();
		return () => window.removeEventListener('resize', updateSize);
	}, []);
	return size;
}

function ItemPage() {
	const queryClient = useQueryClient();
	const { itemId } = useParams();
	const itemQuery = useQuery(ReactQueryUtil.itemKey(itemId), () => Client.getItem(itemId));
	const tagsQuery = useQuery(ReactQueryUtil.TAGS_KEY, Client.getTags);
	let [addTagMode, setAddTagMode] = useState(false);
	let [windowWidth, windowHeight] = useWindowSize();

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
			queryClient.refetchQueries({ queryKey: itemQuery.videoElequeryKey });
		});
	};

	const onTagClicked = (tag) => {
		console.log('send to gallery with selected filter');
	};

	const calcHeight = () => {
		let result = (windowHeight / 10) * 7;

		if (result < 400) {
			return 400;
		}

		return result;
	};

	const calcWidth = () => {
		let ratio = itemQuery.data.width / itemQuery.data.height;
		let actualWidth = calcHeight() * ratio;
		return actualWidth;
	};

	return (
		<Box
			sx={{
				display: 'flex',
				flexDirection: 'column',
				alignItems: 'center',
				padding: '30px 50px',
			}}
		>
			{itemQuery.isSuccess && (
				<Stack flexGrow={1} flexDirection="column" gap="20px" height={calcHeight()} width={calcWidth()}>
					<Player url={itemQuery.data.url} />
					<Typography variant="h5">{itemQuery.data.title}</Typography>
					<Stack flexDirection="row" gap="10px">
						<TagChips
							flexDirection="column"
							tags={itemQuery.data.tags.filter((cur) => TagsUtil.isDirectoriesCategory(cur.parentId))}
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
		</Box>
	);
}

export default ItemPage;
