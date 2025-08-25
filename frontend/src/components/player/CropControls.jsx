import { useTheme } from '@emotion/react';
import CancelIcon from '@mui/icons-material/Cancel';
import CropIcon from '@mui/icons-material/Crop';
import DoneIcon from '@mui/icons-material/Done';
import { IconButton, Tooltip } from '@mui/material';
import Client from '../../utils/client';
import { usePlayerActionStore } from './PlayerActionStore';
import { usePlayerStore } from './PlayerStore';

function CropControls() {
	const theme = useTheme();
	const playerStore = usePlayerStore();
	const playerActionStore = usePlayerActionStore();

	const cropFrame = (second, crop) => {
		Client.cropFrame(playerStore.itemId, second, crop);
	};

	return playerActionStore.cropActive() ? (
		<>
			<IconButton
				onClick={() => {
					let frame = playerActionStore.cropCompleted();
					cropFrame(playerStore.currentTime, frame);
				}}
			>
				<DoneIcon sx={{ fontSize: theme.iconSize(1) }} />
			</IconButton>
			<IconButton onClick={playerActionStore.cropCanceled}>
				<CancelIcon sx={{ fontSize: theme.iconSize(1) }} />
			</IconButton>
		</>
	) : (
		<Tooltip title="Crop">
			<IconButton
				onClick={() => {
					playerStore.pause();
					playerActionStore.startCrop();
				}}
			>
				<CropIcon sx={{ fontSize: theme.iconSize(1) }} />
			</IconButton>
		</Tooltip>
	);
}

export default CropControls;
