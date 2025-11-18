import DownloadIcon from '@mui/icons-material/Download';
import { IconButton, Tooltip, useTheme } from '@mui/material';
import { useRef } from 'react';
import { useSubtitleStore } from './SubtitlesStore';

function SubtitleDownloadButton({ subtitle }) {
	const theme = useTheme();
	const subtitleStore = useSubtitleStore();
	const timeoutRef = useRef(null);

	const clicked = (e) => {
		e.preventDefault();
		e.stopPropagation();
		console.log('download clicked');
	};

	return (
		<Tooltip title="Download">
			<IconButton onClick={clicked} sx={{ color: theme.palette.secondary.main }}>
				<DownloadIcon />
			</IconButton>
		</Tooltip>
	);
}

export default SubtitleDownloadButton;
