import SubtitlesIcon from '@mui/icons-material/Subtitles';
import { IconButton, useTheme } from '@mui/material';
import { useSubtitleStore } from './SubtitlesStore';

function SubtitlesButton() {
	const theme = useTheme();
	const subtitleStore = useSubtitleStore();

	return (
		<IconButton onClick={() => subtitleStore.toggleSubtitlesControls()}>
			<SubtitlesIcon sx={{ fontSize: theme.iconSize(1) }} />
		</IconButton>
	);
}

export default SubtitlesButton;
