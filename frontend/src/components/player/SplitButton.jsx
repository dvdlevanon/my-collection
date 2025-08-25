import { useTheme } from '@emotion/react';
import SplitIcon from '@mui/icons-material/ContentCut';
import { IconButton } from '@mui/material';
import { useQueryClient } from '@tanstack/react-query';
import { useState } from 'react';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import ConfirmationDialog from '../dialogs/ConfirmationDialog';
import { usePlayerStore } from './PlayerStore';

function SplitButton() {
	const theme = useTheme();
	const queryClient = useQueryClient();
	const { itemId, allowToSplit, currentTime } = usePlayerStore();
	const [showSplitVideoConfirmationDialog, setShowSplitVideoConfirmationDialog] = useState(false);
	const [splitVideoSecond, setSplitVideoSecond] = useState(0);

	const closeSplitVideoDialog = () => {
		setSplitVideoSecond(0);
		setShowSplitVideoConfirmationDialog(false);
	};

	const splitItem = () => {
		Client.splitItem(itemId, splitVideoSecond).then(() => {
			ReactQueryUtil.updateItem(queryClient, itemId, true);
		});

		closeSplitVideoDialog();
	};

	const splitClicked = () => {
		setShowSplitVideoConfirmationDialog(true);
		setSplitVideoSecond(currentTime);
	};

	return (
		<>
			<IconButton disabled={!allowToSplit} onClick={splitClicked}>
				<SplitIcon sx={{ fontSize: theme.iconSize(1) }} />
			</IconButton>
			<ConfirmationDialog
				open={showSplitVideoConfirmationDialog}
				title="Split Video"
				text={'Are you sure you want to split the video at second ' + splitVideoSecond + '?'}
				actionButtonTitle="Split"
				onCancel={closeSplitVideoDialog}
				onConfirm={splitItem}
			/>
		</>
	);
}

export default SplitButton;
