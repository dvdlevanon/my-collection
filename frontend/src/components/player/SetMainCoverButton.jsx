import { useTheme } from '@emotion/react';
import ImageIcon from '@mui/icons-material/Image';
import { IconButton } from '@mui/material';
import { useQueryClient } from '@tanstack/react-query';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import { usePlayerStore } from './PlayerStore';

function SetMainCoverButton() {
	const theme = useTheme();
	const queryClient = useQueryClient();
	const playerStore = usePlayerStore();

	const setMainCover = (second) => {
		Client.setMainCover(playerStore.itemId, second).then(() => {
			ReactQueryUtil.updateItem(queryClient, playerStore.itemId, true);
		});
	};

	return (
		<IconButton onClick={() => setMainCover(playerStore.currentTime)}>
			<ImageIcon sx={{ fontSize: theme.iconSize(1) }} />
		</IconButton>
	);
}

export default SetMainCoverButton;
