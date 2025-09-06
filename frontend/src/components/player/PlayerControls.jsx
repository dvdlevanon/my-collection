import { useTheme } from '@emotion/react';
import { Box, Fade, Stack } from '@mui/material';
import AutoPlayButton from './AutoPlayButton';
import CropControls from './CropControls';
import FullscreenButton from './FullscreenButton';
import HighlightButton from './HighlightButton';
import PlayButton from './PlayButton';
import { usePlayerActionStore } from './PlayerActionStore';
import PlayerScrubber from './PlayerScrubber';
import { usePlayerStore } from './PlayerStore';
import SetMainCoverButton from './SetMainCoverButton';
import SplitButton from './SplitButton';
import SubtitlesButton from './subtitles/SubtitlesButton';
import TimeDisplay from './TimeDisplay';
import TimingControls from './TimingControls';
import VolumeControls from './VolumeControls';

function PlayerControls() {
	const playerStore = usePlayerStore();
	const playerActionStore = usePlayerActionStore();
	const theme = useTheme();

	return (
		<Fade in={playerStore.controlsVisible}>
			<Stack
				onMouseEnter={() => playerStore.showControls(false)}
				zIndex={2}
				sx={{
					position: 'absolute',
					padding: theme.spacing(1),
					bottom: -1,
					left: 0,
					right: 0,
					flexDirection: 'column',
					background:
						'linear-gradient(180deg, ' +
						theme.palette.gradient.color1 +
						'00 0%, ' +
						theme.palette.gradient.color2 +
						'FF 100%)',
				}}
			>
				<PlayerScrubber min={playerStore.startTime} max={playerStore.endTime} value={playerStore.currentTime} />
				<Stack flexDirection="row" alignItems="center" gap={theme.spacing(2)}>
					<PlayButton />
					<AutoPlayButton />
					<VolumeControls />
					<TimeDisplay />
					<Box display="flex" flexGrow={1} justifyContent="flex-end">
						{!playerActionStore.cropActive() && (
							<>
								<TimingControls setRelativeTime={playerStore.offsetSeek} />
								<SetMainCoverButton />
								<SubtitlesButton />
							</>
						)}
						<CropControls />
						{!playerActionStore.cropActive() && (
							<>
								<SplitButton />
								<HighlightButton />
								<FullscreenButton />
							</>
						)}
					</Box>
				</Stack>
			</Stack>
		</Fade>
	);
}

export default PlayerControls;
