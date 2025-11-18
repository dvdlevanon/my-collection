import { Stack, useTheme } from '@mui/material';
import AvailableSubtitlesChooser from './AvailableSubtitlesChooser';
import SubtitileSyncer from './SubtitleSyncer';

function SubtitlesSettings() {
	const theme = useTheme();

	return (
		<Stack gap={theme.spacing(1)}>
			<AvailableSubtitlesChooser />
			<SubtitileSyncer />
		</Stack>
	);
}

export default SubtitlesSettings;
