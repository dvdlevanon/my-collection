import { Stack, useTheme } from '@mui/material';
import AvailableSubtitlesChooser from './AvailableSubtitlesChooser';
import { useSubtitleStore } from './SubtitlesStore';
import SubtitileSyncer from './SubtitleSyncer';

function SubtitlesSettings() {
	const theme = useTheme();
	const subtitleStore = useSubtitleStore();

	return (
		<Stack gap={theme.spacing(1)}>
			<AvailableSubtitlesChooser />
			<SubtitileSyncer />
		</Stack>
	);
}

export default SubtitlesSettings;
