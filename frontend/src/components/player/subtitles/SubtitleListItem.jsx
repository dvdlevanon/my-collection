import LocalIcon from '@mui/icons-material/Description';
import OnlineIcon from '@mui/icons-material/Language';
import { Box, ListItemButton, ListItemIcon, ListItemText, useTheme } from '@mui/material';
import SubtitleDeleteButton from './SubtitleDeleteButton';
import SubtitleDownloadButton from './SubtitleDownloadButton';
import { useSubtitleStore } from './SubtitlesStore';

function SubtitleListItem({ subtitle, refetchOnlineSubtitles }) {
	const theme = useTheme();
	const subtitleStore = useSubtitleStore();

	const isReady = () => {
		return subtitle.id === 'local' || subtitle.url !== '';
	};

	const isOnlineAndReady = () => {
		return subtitle.id !== 'local' && subtitle.url !== '';
	};

	const isSelected = () => {
		return subtitle.url != '' && subtitleStore.selectedSubtitleUrl == subtitle.url;
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
				border: isSelected() ? theme.border(1, 'solid', theme.palette.primary.main) : 'auto',
			}}
		>
			<ListItemIcon>{subtitle.id == 'local' ? <LocalIcon /> : <OnlineIcon />}</ListItemIcon>
			<ListItemText sx={{ color: isReady() ? 'auto' : theme.palette.primary.disabled }}>
				{subtitle.title}
			</ListItemText>
			{!isReady() && (
				<Box sx={{ pointerEvents: 'auto' }}>
					<SubtitleDownloadButton subtitle={subtitle} refetchOnlineSubtitles={refetchOnlineSubtitles} />
				</Box>
			)}
			{isOnlineAndReady() && (
				<Box sx={{ pointerEvents: 'auto' }}>
					<SubtitleDeleteButton subtitle={subtitle} refetchOnlineSubtitles={refetchOnlineSubtitles} />
				</Box>
			)}
		</ListItemButton>
	);
}

export default SubtitleListItem;
