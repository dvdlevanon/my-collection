import { Stack, useTheme } from '@mui/material';
import SubtitlesFilter from './SubtitlesFilter';
import { useSubtitleStore } from './SubtitlesStore';
import SubtitileSyncer from './SubtitleSyncer';

function SubtitlesFinder() {
	const theme = useTheme();
	const subtitleStore = useSubtitleStore();

	return (
		<Stack gap={theme.spacing(1)}>
			<SubtitlesFilter />
			<SubtitileSyncer />
		</Stack>
	);
}

export default SubtitlesFinder;
