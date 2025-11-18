import { CircularProgress, List, ListItem, ListItemButton, ListItemText, Stack, useTheme } from '@mui/material';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { useEffect, useRef } from 'react';
import ReactQueryUtil from '../../../utils/react-query-util';
import { usePlayerStore } from '../PlayerStore';
import SubtitleListItem from './SubtitleListItem';

function AvailableSubtitlesChooser() {
	const theme = useTheme();
	const containerRef = useRef();
	const playerStore = usePlayerStore();
	const availableSubtitleQuery = useQuery(ReactQueryUtil.availableSubtitleQuery(playerStore.itemId));
	const onlineSubtitleQuery = useQuery(ReactQueryUtil.onlineSubtitleQuery(playerStore.itemId, 'he', false));
	const queryClient = useQueryClient();

	useEffect(() => {
		if (availableSubtitleQuery.data) {
		}
	}, [availableSubtitleQuery.data]);

	const refetchOnlineSubtitles = () => {
		queryClient.invalidateQueries({
			queryKey: ReactQueryUtil.onlineSubtitleQueryKey(playerStore.itemId, 'he', false),
		});
	};

	const getSubtitles = () => {
		if (!availableSubtitleQuery.data) {
			return [];
		}

		return availableSubtitleQuery.data;
	};

	const getOnlineSubtitles = () => {
		if (!onlineSubtitleQuery.data) {
			return [];
		}

		return onlineSubtitleQuery.data;
	};

	return (
		<Stack ref={containerRef} maxHeight={'30vh'} overflow={'auto'}>
			<List>
				{getSubtitles().map((subtitle) => {
					return (
						<ListItem key={subtitle.url} value={subtitle.url}>
							<SubtitleListItem subtitle={subtitle} />
						</ListItem>
					);
				})}
				{!onlineSubtitleQuery.data && (
					<ListItem key={'loading'} value={'loading'}>
						<ListItemButton disabled sx={{ gap: theme.spacing(1) }}>
							<CircularProgress />
							<ListItemText sx={{ color: theme.palette.primary.disabled }}>
								Loading Online Subtitles...
							</ListItemText>
						</ListItemButton>
					</ListItem>
				)}
				{getOnlineSubtitles().map((subtitle) => {
					return (
						<ListItem key={subtitle.id} value={subtitle.url}>
							<SubtitleListItem subtitle={subtitle} refetchOnlineSubtitles={refetchOnlineSubtitles} />
						</ListItem>
					);
				})}
			</List>
		</Stack>
	);
}

export default AvailableSubtitlesChooser;
