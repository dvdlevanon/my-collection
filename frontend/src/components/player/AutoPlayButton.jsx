import { useTheme } from '@emotion/react';
import PlayNextIcon from '@mui/icons-material/QueuePlayNext';
import { IconButton, Tooltip } from '@mui/material';
import { usePlayerStore } from './PlayerStore';

function AutoPlayButton() {
	const theme = useTheme();
	const playerStore = usePlayerStore();

	return (
		<Tooltip title={'Toggle Auto Play Next'}>
			<IconButton
				onClick={() => {
					playerStore.setAutoPlayNext(!playerStore.autoPlayNext);
				}}
			>
				{
					<PlayNextIcon
						color={playerStore.autoPlayNext ? 'secondary' : 'auto'}
						sx={{ fontSize: theme.iconSize(1) }}
					/>
				}
			</IconButton>
		</Tooltip>
	);
}

export default AutoPlayButton;
