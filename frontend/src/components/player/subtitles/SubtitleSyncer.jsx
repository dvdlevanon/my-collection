import { Slider, Stack, Typography, useTheme } from '@mui/material';
import { useEffect, useRef } from 'react';
import { useSubtitleStore } from './SubtitlesStore';

function SubtitileSyncer() {
	const theme = useTheme();
	const subtitleStore = useSubtitleStore();
	const timeoutRef = useRef(null);

	useEffect(() => {
		return () => {
			if (timeoutRef.current) {
				clearTimeout(timeoutRef.current);
			}
		};
	}, []);

	const onChange = (newValue) => {
		if (timeoutRef.current) {
			clearTimeout(timeoutRef.current);
		}

		subtitleStore.setSubtitleOffsetMillis(newValue * 1000);
		timeoutRef.current = setTimeout(() => {
			subtitleStore.setSubtitleOffsetMillis(newValue * 1000);
		}, 150);
	};

	return (
		<Stack flexDirection={'row'} alignItems={'center'} gap={theme.spacing(3)} justifyContent={'space-between'}>
			<Typography>Time Offset</Typography>
			<Slider
				min={-60}
				max={60}
				step={0.1}
				value={subtitleStore.subtitleOffsetMillis / 1000}
				size="small"
				valueLabelDisplay="on"
				onChange={(e, newValue) => onChange(newValue)}
				sx={{
					width: '400px',
					padding: theme.multiSpacing(4, 0),
				}}
			/>
		</Stack>
	);
}

export default SubtitileSyncer;
