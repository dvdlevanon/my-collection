import { Box, Typography } from '@mui/material';
import { useQuery } from '@tanstack/react-query';
import { useEffect, useRef, useState } from 'react';
import ReactQueryUtil from '../../utils/react-query-util';
import { usePlayerStore } from './PlayerStore';

function getCurrentSubtitle(currentTimeSeconds, subtitleData, currentIndexRef) {
	if (!subtitleData || !subtitleData.items || subtitleData.items.length === 0) {
		return '';
	}

	const currentTimeMillis = currentTimeSeconds * 1000;
	const items = subtitleData.items;

	if (currentIndexRef.current >= 0 && currentIndexRef.current < items.length) {
		const currentItem = items[currentIndexRef.current];
		if (currentTimeMillis >= currentItem.start_millis && currentTimeMillis <= currentItem.end_millis) {
			return currentItem.text;
		}
	}

	for (let i = 0; i < items.length; i++) {
		const item = items[i];
		if (currentTimeMillis >= item.start_millis && currentTimeMillis <= item.end_millis) {
			currentIndexRef.current = i;
			return item.text;
		}
	}

	currentIndexRef.current = -1;
	return '';
}

function Subtitles({ itemId }) {
	const playerStore = usePlayerStore();
	const [text, setText] = useState('');
	const currentIndexRef = useRef(-1);
	const subtitleQuery = useQuery(ReactQueryUtil.subtitleQuery(itemId));

	useEffect(() => {
		if (subtitleQuery.data) {
			let text = getCurrentSubtitle(playerStore.currentTime, subtitleQuery.data, currentIndexRef);
			setText(text);
		}
	}, [playerStore.currentTime, subtitleQuery.data]);

	return (
		<Box
			position="absolute"
			bottom={playerStore.controlsVisible ? 100 : 20}
			left={0}
			right={0}
			display={'flex'}
			flexDirection={'column'}
			sx={{ pointerEvents: 'none' }}
		>
			<Typography fontSize={60} textAlign={'center'} variant="caption">
				{text}
			</Typography>
		</Box>
	);
}

export default Subtitles;
