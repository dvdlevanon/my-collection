import { Box, Slider, Typography } from '@mui/material';
import { useState } from 'react';
import TimeUtil from '../../utils/time-utils';
import { usePlayerStore } from './PlayerStore';

function PlayerScrubber({ min, max, value }) {
	const [showTime, setShowTime] = useState(false);
	const [mouseX, setMouseX] = useState(0);
	const playerStore = usePlayerStore();

	return (
		<Box
			onMouseEnter={() => setShowTime(false)}
			onMouseLeave={() => setShowTime(false)}
			onMouseMove={(e) => {
				let bounds = e.currentTarget.getBoundingClientRect();
				let x = Math.floor(e.clientX - bounds.left);

				if (x > bounds.width - 100) {
					x = bounds.width - 100;
				}

				setMouseX(x < 0 ? 0 : x);
			}}
			sx={{
				position: 'relative',
			}}
		>
			{showTime && (
				<Typography sx={{ position: 'absolute', left: mouseX, top: -10 }}>
					{TimeUtil.formatSeconds(value - min)} / {TimeUtil.formatSeconds(max - min)} - {mouseX}
				</Typography>
			)}
			<Slider
				min={min}
				max={max}
				value={value}
				valueLabelDisplay="auto"
				valueLabelFormat={(number) => {
					return TimeUtil.formatSeconds(value - min);
				}}
				onChange={(e, newValue) => playerStore.seek(newValue)}
			/>
		</Box>
	);
}

export default PlayerScrubber;
