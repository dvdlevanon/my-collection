import { Slider, Stack, Typography, useTheme } from '@mui/material';

function SubtitileSyncer() {
	const theme = useTheme();

	return (
		<Stack flexDirection={'row'} alignItems={'center'} gap={theme.spacing(3)} justifyContent={'space-between'}>
			<Typography>Time Offset</Typography>
			<Slider
				min={-30}
				max={30}
				step={0.1}
				value={0}
				size="small"
				valueLabelDisplay="on"
				sx={{
					width: '200px',
					padding: theme.multiSpacing(4, 0),
				}}
			/>
		</Stack>
	);
}

export default SubtitileSyncer;
