import { Box, Typography } from '@mui/material';
import { useQuery } from '@tanstack/react-query';
import { useEffect, useRef, useState } from 'react';
import ReactQueryUtil from '../../../utils/react-query-util';
import { usePlayerStore } from '../PlayerStore';
import { useSubtitleStore } from './SubtitlesStore';

function getCurrentSubtitle(currentTimeSeconds, subtitleData, currentIndexRef, subtitleOffset) {
	if (!subtitleData || !subtitleData.items || subtitleData.items.length === 0) {
		return '';
	}

	const currentTimeMillis = currentTimeSeconds * 1000 + subtitleOffset;
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

function Subtitles() {
	const playerStore = usePlayerStore();
	const subtitleStore = useSubtitleStore();
	const [text, setText] = useState('');
	const currentIndexRef = useRef(-1);
	const subtitleQuery = useQuery(ReactQueryUtil.subtitleQuery(subtitleStore.selectedSubtitleUrl));

	useEffect(() => {
		subtitleStore.loadFromLocalStorage();
	}, []);

	useEffect(() => {
		if (subtitleQuery.data) {
			let text = getCurrentSubtitle(
				playerStore.currentTime,
				subtitleQuery.data,
				currentIndexRef,
				subtitleStore.subtitleOffsetMillis
			);
			setText(text);
		}
	}, [playerStore.currentTime, subtitleQuery.data, subtitleStore.subtitleOffsetMillis]);

	const getTextShadaw = () => {
		let shadowWidth = subtitleStore.fontShadowWidth;
		let shadawColor = subtitleStore.fontShadowColor;

		return `
			-${shadowWidth}px -${shadowWidth}px 0 ${shadawColor},
			0 -${shadowWidth}px 0 ${shadawColor},
			${shadowWidth}px -${shadowWidth}px 0 ${shadawColor},
			${shadowWidth}px 0 0 ${shadawColor},
			${shadowWidth}px ${shadowWidth}px 0 ${shadawColor},
			0 ${shadowWidth}px 0 ${shadawColor},
			-${shadowWidth}px ${shadowWidth}px 0 ${shadawColor},
			-${shadowWidth}px 0 0 ${shadawColor};
		`;
	};

	return (
		<Box
			position="absolute"
			bottom={playerStore.controlsVisible ? 100 : 20}
			left={0}
			right={0}
			display={'flex'}
			flexDirection={'column'}
			sx={{ pointerEvents: 'none', containerType: 'inline-size' }}
		>
			<Typography
				fontSize={'clamp(16px, ' + subtitleStore.fontSize + 'cqw, 1000px)'}
				textAlign={'center'}
				variant="caption"
				color={subtitleStore.fontColor}
				sx={{
					textShadow: getTextShadaw(),
				}}
			>
				{text}
			</Typography>
		</Box>
	);
}

export default Subtitles;
