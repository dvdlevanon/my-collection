import DeleteIcon from '@mui/icons-material/Delete';
import { CircularProgress, IconButton, Tooltip, useTheme } from '@mui/material';
import { useState } from 'react';
import Client from '../../../utils/client';
import { usePlayerStore } from '../PlayerStore';

function SubtitleDeleteButton({ subtitle, refetchOnlineSubtitles }) {
	const theme = useTheme();
	const playerStore = usePlayerStore();
	const [isDeleting, setIsDeleting] = useState(false);

	const clicked = (e) => {
		e.preventDefault();
		e.stopPropagation();
		if (isDeleting) {
			return;
		}

		setIsDeleting(true);
		Client.deleteSubtitle(subtitle.url).then(() => {
			setIsDeleting(false);
			refetchOnlineSubtitles();
		});
	};

	return (
		<Tooltip title="Delete">
			<IconButton onClick={clicked} sx={{ color: theme.palette.secondary.main }}>
				{isDeleting ? <CircularProgress size={theme.iconSize(1)} /> : <DeleteIcon size={theme.iconSize(1)} />}
			</IconButton>
		</Tooltip>
	);
}

export default SubtitleDeleteButton;
