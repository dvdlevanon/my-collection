import CloseIcon from '@mui/icons-material/Close';
import { IconButton, Stack, Typography, useTheme } from '@mui/material';
import { useSubtitleStore } from './SubtitlesStore';

function SubtitlesControlsHeader() {
	const theme = useTheme();
	const subtitleStore = useSubtitleStore();

	return (
		<Stack flexDirection={'row'} alignItems={'center'} gap={theme.spacing(1)}>
			<IconButton
				onClick={(e) => {
					e.preventDefault();
					e.stopPropagation();
					subtitleStore.hideSubtitlesControls();
				}}
			>
				<CloseIcon sx={{ fontSize: theme.iconSize(1) }} />
			</IconButton>
			<Typography
				variant="body1"
				noWrap
				onClick={(e) => {
					e.preventDefault();
					e.stopPropagation();
				}}
			>
				Subtitles Options
			</Typography>
		</Stack>
	);
}

export default SubtitlesControlsHeader;
