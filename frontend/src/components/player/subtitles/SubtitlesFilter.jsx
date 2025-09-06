import { Checkbox, FormControlLabel, Stack, TextField } from '@mui/material';
import TimeUtil from '../../../utils/time-utils';
import { usePlayerStore } from '../PlayerStore';

function SubtitlesFilter() {
	const playerStore = usePlayerStore();

	return (
		<Stack>
			<FormControlLabel
				label={'Must have subtitles now ' + TimeUtil.formatDuration(playerStore.currentTime)}
				control={<Checkbox />}
			/>
			<TextField placeholder="Subtitle must contain" />
		</Stack>
	);
}

export default SubtitlesFilter;
