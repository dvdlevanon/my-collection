import DeleteIcon from '@mui/icons-material/Delete';
import { Box, IconButton, Stack } from '@mui/material';
import React, { useState } from 'react';
import AspectRatioUtil from '../../utils/aspect-ratio-util';
import ConfirmationDialog from '../dialogs/ConfirmationDialog';
import Item from '../items-viewer/Item';

function SubItem({ item, itemWidth, highlighted, onDeleteItem }) {
	const [showDelete, setShowDelete] = useState(false);
	const [showConfirmDialog, setShowConfirmDialog] = useState(false);

	return (
		<Stack
			flexDirection="row"
			gap="10px"
			padding="10px"
			sx={{
				cursor: 'pointer',
				borderRadius: '10px',
				backgroundColor: highlighted ? 'dark.lighter2' : 'unset',
				'&:hover': {
					backgroundColor: 'dark.lighter',
				},
			}}
		>
			<Box position="relative" onMouseEnter={() => setShowDelete(true)} onMouseLeave={() => setShowDelete(false)}>
				<Item
					item={item}
					preferPreview={true}
					itemWidth={itemWidth}
					itemHeight={AspectRatioUtil.calcHeight(itemWidth, AspectRatioUtil.asepctRatio16_9)}
					direction="row"
					showOffests={true}
					titleSx={{
						whiteSpace: 'normal',
						lineHeight: '1.5em',
						maxHeight: '3em',
						textAlign: 'start',
					}}
					itemLinkBuilder={(item) => {
						return '/spa/item/' + item.id;
					}}
				/>
				{showDelete && onDeleteItem && (
					<IconButton
						onClick={() => {
							setShowConfirmDialog(true);
							setShowDelete(false);
						}}
						sx={{
							position: 'absolute',
							bottom: '5px',
							right: '5px',
						}}
					>
						<DeleteIcon />
					</IconButton>
				)}
				{showConfirmDialog && (
					<ConfirmationDialog
						title="Delete Item"
						text={'Are you sure you want to delete ' + item.title}
						actionButtonTitle="Delete"
						onCancel={() => {
							setShowConfirmDialog(false);
							setShowDelete(false);
						}}
						onConfirm={() => {
							setShowConfirmDialog(false);
							setShowDelete(false);
							onDeleteItem(item);
						}}
					/>
				)}
			</Box>
		</Stack>
	);
}

export default SubItem;
