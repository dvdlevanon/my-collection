import AddIcon from '@mui/icons-material/Add';
import DeleteIcon from '@mui/icons-material/Delete';
import { Box, Chip, IconButton, Stack, Tooltip, Typography } from '@mui/material';
import { useLayoutEffect, useState } from 'react';
import { useQuery, useQueryClient } from 'react-query';
import { useNavigate, useParams } from 'react-router-dom';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import TagsUtil from '../../utils/tags-util';
import AttachTagDialog from '../dialogs/AttachTagDialog';
import ConfirmationDialog from '../dialogs/ConfirmationDialog';
import Highlights from '../highlights/Highlights';
import ItemTitle from '../items-viewer/ItemTitle';
import Player from '../player/Player';
import SubItems from '../sub-items/SubItems';
import TagChips from '../tags-chip/TagChips';

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
	const suggestedQuery = useQuery({
		queryKey: ReactQueryUtil.suggestedItemsKey(itemId),
		queryFn: () => Client.getSuggestedItems(itemId),
		staleTime: Infinity,
		cacheTime: Infinity,
	});
	const [addTagMode, setAddTagMode] = useState(false);
	const [windowWidth, windowHeight] = useWindowSize();
	const [showSplitVideoConfirmationDialog, setShowSplitVideoConfirmationDialog] = useState(false);
	const [showDeleteItemConfirmationDialog, setShowDeleteItemConfirmationDialog] = useState(false);
	const [splitVideoSecond, setSplitVideoSecond] = useState(0);
	const itemQuery = useQuery({
		queryKey: ReactQueryUtil.itemKey(itemId),
		queryFn: () => Client.getItem(itemId),
		onSuccess: (item) => {
			document.title = item.title;
		},
	});
	const navigate = useNavigate();

	const onAddTag = () => {
		setAddTagMode(true);
	};

	const onTitleChanged = (newTitle) => {
		if (!newTitle) {
			return;
		}

		Client.saveItem({ ...itemQuery.data, title: newTitle }, () => {
			queryClient.refetchQueries({ queryKey: itemQuery.queryKey });
		});
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

	const setMainCover = (second) => {
		Client.setMainCover(itemQuery.data.id, second).then(() => {
			queryClient.refetchQueries({ queryKey: itemQuery.queryKey });
		});
	};

	const onDeleteItem = (item, deleteRealFile) => {
		Client.deleteItem(item.id, deleteRealFile).then(() => {
			if (deleteRealFile) {
				navigate('/');
			} else {
				queryClient.refetchQueries({ queryKey: itemQuery.queryKey });
			}
		});
	};

	const closeSplitVideoDialog = () => {
		setSplitVideoSecond(0);
		setShowSplitVideoConfirmationDialog(false);
	};

	const splitItem = () => {
		Client.splitItem(itemQuery.data.id, splitVideoSecond).then(() => {
			queryClient.refetchQueries({ queryKey: itemQuery.queryKey });
		});
		closeSplitVideoDialog();
	};

	const makeHighlight = (startSecond, endSecond, highlightId) => {
		Client.makeHighlight(itemQuery.data.id, startSecond, endSecond, highlightId).then(() => {
			queryClient.refetchQueries({ queryKey: itemQuery.queryKey });
		});
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

	const shouldShowSubItems = () => {
		return itemQuery.data.main_item || itemQuery.data.sub_items;
	};

	const shouldShowHighlights = () => {
		return itemQuery.data.highlight_parent_id || itemQuery.data.highlights;
	};

	return (
		<Box
			sx={{
				display: 'flex',
				flexDirection: 'row',
				padding: '30px 50px',
				justifyContent: 'center',
				gap: '10px',
			}}
		>
			<Stack maxWidth={500}>
				{itemQuery.isSuccess && shouldShowHighlights() && <Highlights item={itemQuery.data} />}
			</Stack>
			<Box
				sx={{
					display: 'flex',
					flexDirection: 'column',
					alignItems: 'center',
				}}
			>
				{itemQuery.isSuccess && suggestedQuery.isSuccess && (
					<Stack flexGrow={1} flexDirection="column" gap="20px" height={calcHeight()} width={calcWidth()}>
						<Player
							url={itemQuery.data.url}
							suggestedItems={suggestedQuery.data}
							setMainCover={setMainCover}
							allowToSplit={() => {
								return !itemQuery.data.sub_items;
							}}
							splitVideo={(splitSecond) => {
								setShowSplitVideoConfirmationDialog(true);
								setSplitVideoSecond(splitSecond);
							}}
							makeHighlight={makeHighlight}
							startPosition={itemQuery.data.start_position || 0}
							initialEndPosition={itemQuery.data.end_position || 0}
						/>
						<ItemTitle
							item={itemQuery.data}
							variant="h5"
							onTagAdded={onTagAdded}
							onTitleChanged={onTitleChanged}
							sx={{
								whiteSpace: 'normal',
								overflow: 'visible',
								textAlign: 'start',
							}}
							withTooltip={false}
							withMenu={true}
						/>
						<Stack flexDirection="row" gap="10px" alignItems="center">
							<TagChips
								flexDirection="column"
								tags={(itemQuery.data.tags || []).filter((cur) =>
									TagsUtil.allowToAddToCategory(cur.parentId)
								)}
								linkable={true}
								onDelete={onTagRemoved}
								onClick={() => {}}
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
							<Stack width="100%" alignItems="flex-end">
								<Tooltip title="Delete this item">
									<IconButton onClick={() => setShowDeleteItemConfirmationDialog(true)}>
										<DeleteIcon />
									</IconButton>
								</Tooltip>
							</Stack>
						</Stack>
						<Typography variant="body2" color="bright.darker2" padding="0px 10px">
							Resolution: {itemQuery.data.width} * {itemQuery.data.height}
							<br />
							Video Codec: {itemQuery.data.video_codec}
							<br />
							Audio Codec: {itemQuery.data.audio_codec}
						</Typography>
						<AttachTagDialog
							open={addTagMode}
							item={itemQuery.data}
							onTagAdded={onTagAdded}
							onClose={(e) => setAddTagMode(false)}
						/>
					</Stack>
				)}
			</Box>
			<Stack maxWidth={500}>
				{itemQuery.isSuccess && shouldShowSubItems() && (
					<SubItems item={itemQuery.data} onDeleteItem={(item) => onDeleteItem(item, false)} />
				)}
			</Stack>
			{showSplitVideoConfirmationDialog && (
				<ConfirmationDialog
					title="Split Video"
					text={'Are you sure you want to split the video at second ' + splitVideoSecond + '?'}
					actionButtonTitle="Split"
					onCancel={closeSplitVideoDialog}
					onConfirm={splitItem}
				/>
			)}
			{showDeleteItemConfirmationDialog && (
				<ConfirmationDialog
					title="Delete Item"
					text={'Are you sure you want to delete ' + itemQuery.data.title + '?'}
					actionButtonTitle="Delete"
					onCancel={() => setShowDeleteItemConfirmationDialog(false)}
					onConfirm={() => onDeleteItem(itemQuery.data, true)}
				/>
			)}
		</Box>
	);
}

export default ItemPage;
