import { useTheme } from '@emotion/react';
import DeleteIcon from '@mui/icons-material/Delete';
import { Box, IconButton, Stack } from '@mui/material';
import React, { useState } from 'react';
import AspectRatioUtil from '../../utils/aspect-ratio-util';
import ConfirmationDialog from '../dialogs/ConfirmationDialog';
import Item from '../items-viewer/Item';

function SubItem({ item, itemWidth, highlighted, onDeleteItem }) {
	const [showDelete, setShowDelete] = useState(false);
	const [showConfirmDialog, setShowConfirmDialog] = useState(false);
	const theme = useTheme();

	return (
		<Stack
			flexDirection="row"
			gap={theme.spacing(1)}
			padding={theme.spacing(1)}
			sx={{
				cursor: 'pointer',
				borderRadius: theme.spacing(1),
				backgroundColor: highlighted ? theme.palette.primary.light : 'unset',
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
							bottom: theme.spacing(0.5),
							right: theme.spacing(0.5),
						}}
					>
						<DeleteIcon sx={{ fontSize: theme.iconSize(1) }} />
					</IconButton>
				)}
				<ConfirmationDialog
					open={showConfirmDialog}
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
			</Box>
		</Stack>
	);
}

export default SubItem;
