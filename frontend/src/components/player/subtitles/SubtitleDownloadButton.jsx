import DownloadIcon from '@mui/icons-material/Download';
import { CircularProgress, IconButton, Tooltip, useTheme } from '@mui/material';
import { useState } from 'react';
import Client from '../../../utils/client';
import { usePlayerStore } from '../PlayerStore';

function SubtitleDownloadButton({ subtitle, refetchOnlineSubtitles }) {
	const theme = useTheme();
	const playerStore = usePlayerStore();
	const [isDownloading, setIsDownloading] = useState(false);

	const clicked = (e) => {
		e.preventDefault();
		e.stopPropagation();
		if (isDownloading) {
			return;
		}

		setIsDownloading(true);
		Client.downloadSubtitle(playerStore.itemId, subtitle.id, subtitle.title).then(() => {
			setIsDownloading(false);
			refetchOnlineSubtitles();
		});
	};

	return (
		<Tooltip title="Download">
			<IconButton onClick={clicked} sx={{ color: theme.palette.secondary.main }}>
				{isDownloading ? (
					<CircularProgress size={theme.iconSize(1)} />
				) : (
					<DownloadIcon size={theme.iconSize(1)} />
				)}
			</IconButton>
		</Tooltip>
	);
}

export default SubtitleDownloadButton;
