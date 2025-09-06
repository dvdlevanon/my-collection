import { Slider, Stack, Typography, useTheme } from '@mui/material';
import ColorPicker from '../../color-picker/ColorPicker';
import { useSubtitleStore } from './SubtitlesStore';

function SubtitlesAppearanceControls() {
	const theme = useTheme();
	const subtitleStore = useSubtitleStore();

	return (
		<Stack>
			<Stack flexDirection={'row'} alignItems={'center'} gap={theme.spacing(3)} justifyContent={'space-between'}>
				<Typography>Color</Typography>
				<ColorPicker color={subtitleStore.fontColor} onChange={(c) => subtitleStore.setFontColor(c)} />
			</Stack>
			<Stack flexDirection={'row'} alignItems={'center'} gap={theme.spacing(3)} justifyContent={'space-between'}>
				<Typography>Shadow Color</Typography>
				<ColorPicker
					color={subtitleStore.fontShadowColor}
					onChange={(c) => subtitleStore.setFontShadowColor(c)}
				/>
			</Stack>
			<Stack flexDirection={'row'} alignItems={'center'} gap={theme.spacing(3)} justifyContent={'space-between'}>
				<Typography>Font Size</Typography>
				<Slider
					size="small"
					step={0.1}
					min={1}
					max={10}
					value={subtitleStore.fontSize}
					onChange={(e, newValue) => subtitleStore.setFontSize(newValue)}
					valueLabelDisplay="on"
					sx={{
						width: '200px',
						padding: theme.multiSpacing(4, 0),
					}}
				></Slider>
			</Stack>
			<Stack flexDirection={'row'} alignItems={'center'} gap={theme.spacing(3)} justifyContent={'space-between'}>
				<Typography>Shadow Size</Typography>
				<Slider
					size="small"
					step={1}
					min={0}
					max={10}
					value={subtitleStore.fontShadowWidth}
					onChange={(e, newValue) => subtitleStore.setFontShadowWidth(newValue)}
					valueLabelDisplay="on"
					sx={{
						width: '200px',
						padding: theme.multiSpacing(4, 0),
					}}
				></Slider>
			</Stack>
		</Stack>
	);
}

export default SubtitlesAppearanceControls;
