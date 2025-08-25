import { Box, Typography } from '@mui/material';
import TimeUtil from '../../utils/time-utils';
import { usePlayerStore } from './PlayerStore';

function TimeDisplay() {
	const playerStore = usePlayerStore();

	return (
		<Box>
			<Typography>
				{TimeUtil.formatSeconds(playerStore.currentTime - playerStore.startTime)} /{' '}
				{TimeUtil.formatSeconds(playerStore.duration)}
			</Typography>
		</Box>
	);
}

export default TimeDisplay;
