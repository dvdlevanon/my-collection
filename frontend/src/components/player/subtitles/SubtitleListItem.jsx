import LocalIcon from '@mui/icons-material/Description';
import OnlineIcon from '@mui/icons-material/Language';
import { Box, ListItemButton, ListItemIcon, ListItemText, useTheme } from '@mui/material';
import SubtitleDownloadButton from './SubtitleDownloadButton';
import { useSubtitleStore } from './SubtitlesStore';

function SubtitleListItem({ subtitle }) {
	const theme = useTheme();
	const subtitleStore = useSubtitleStore();

	const isReady = () => {
		return subtitle.id === 'local' || subtitle.url !== '';
	};

	const clicked = () => {
		if (isReady()) {
			subtitleStore.setSelectedSubtitleUrl(subtitle.url);
		}
	};

	return (
		<ListItemButton
			onClick={clicked}
			sx={{
				pointerEvents: !isReady() ? 'none' : 'auto',
				gap: theme.spacing(1),
			}}
		>
			<ListItemIcon>{subtitle.id == 'local' ? <LocalIcon /> : <OnlineIcon />}</ListItemIcon>
			<ListItemText sx={{ color: isReady() ? 'auto' : theme.palette.primary.disabled }}>
				{subtitle.title}
			</ListItemText>
			{!isReady() && (
				<Box sx={{ pointerEvents: 'auto' }}>
					<SubtitleDownloadButton />
				</Box>
			)}
		</ListItemButton>
	);
}

export default SubtitleListItem;
