import { useTheme } from '@emotion/react';
import CancelIcon from '@mui/icons-material/Cancel';
import SplitIcon from '@mui/icons-material/ContentCut';
import CropIcon from '@mui/icons-material/Crop';
import DoneIcon from '@mui/icons-material/Done';
import FullscreenIcon from '@mui/icons-material/Fullscreen';
import FullscreenExitIcon from '@mui/icons-material/FullscreenExit';
import ImageIcon from '@mui/icons-material/Image';
import PauseIcon from '@mui/icons-material/Pause';
import PlayIcon from '@mui/icons-material/PlayArrow';
import PlayNextIcon from '@mui/icons-material/QueuePlayNext';
import HighlightIcon from '@mui/icons-material/Stars';
import { Box, Fade, IconButton, Stack, Tooltip, Typography } from '@mui/material';
import TimeUtil from '../../utils/time-utils';
import { usePlayerActionStore } from './PlayerActionStore';
import PlayerSlider from './PlayerSlider';
import { usePlayerStore } from './PlayerStore';
import TimingControls from './TimingControls';
import VolumeControls from './VolumeControls';

function PlayerControls({ setMainCover, splitVideo, allowToSplit, cropFrame }) {
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
				<PlayerSlider
					min={playerStore.startTime}
					max={playerStore.endTime}
					value={playerStore.currentTime}
					onChange={(e, newValue) => playerStore.seek(newValue)}
				/>
				<Stack flexDirection="row" alignItems="center" gap={theme.spacing(2)}>
					<Tooltip title={playerStore.isPlaying ? 'Pause' : 'Play'}>
						<IconButton onClick={playerStore.togglePlay}>
							{playerStore.isPlaying ? (
								<PauseIcon sx={{ fontSize: theme.iconSize(1) }} />
							) : (
								<PlayIcon sx={{ fontSize: theme.iconSize(1) }} />
							)}
						</IconButton>
					</Tooltip>
					<Tooltip title={'Toggle Auto Play Next'}>
						<IconButton
							onClick={() => {
								playerStore.setAutoPlayNext(!playerStore.autoPlayNext);
							}}
						>
							{
								<PlayNextIcon
									color={playerStore.autoPlayNext ? 'secondary' : 'auto'}
									sx={{ fontSize: theme.iconSize(1) }}
								/>
							}
						</IconButton>
					</Tooltip>
					<VolumeControls />
					<Box>
						<Typography>
							{TimeUtil.formatSeconds(playerStore.currentTime - playerStore.startTime)} /{' '}
							{TimeUtil.formatSeconds(playerStore.duration)}
						</Typography>
					</Box>
					<Box display="flex" flexGrow={1} justifyContent="flex-end">
						<TimingControls setRelativeTime={playerStore.offsetSeek} />
						{!playerActionStore.cropActive() && (
							<IconButton onClick={() => setMainCover(playerStore.currentTime)}>
								<ImageIcon sx={{ fontSize: theme.iconSize(1) }} />
							</IconButton>
						)}
						{playerActionStore.cropActive() ? (
							<>
								<IconButton
									onClick={() => {
										let frame = playerActionStore.cropCompleted();
										cropFrame(playerStore.currentTime, frame);
									}}
								>
									<DoneIcon sx={{ fontSize: theme.iconSize(1) }} />
								</IconButton>
								<IconButton onClick={playerActionStore.cropCanceled}>
									<CancelIcon sx={{ fontSize: theme.iconSize(1) }} />
								</IconButton>
							</>
						) : (
							<IconButton
								onClick={() => {
									playerStore.pause();
									playerActionStore.startCrop();
								}}
							>
								<CropIcon sx={{ fontSize: theme.iconSize(1) }} />
							</IconButton>
						)}
						{!playerActionStore.cropActive() && (
							<IconButton disabled={!allowToSplit()} onClick={() => splitVideo(playerStore.currentTime)}>
								<SplitIcon sx={{ fontSize: theme.iconSize(1) }} />
							</IconButton>
						)}
						{!playerActionStore.cropActive() && (
							<IconButton
								onClick={() => playerActionStore.startHighlightCreation(playerStore.currentTime)}
							>
								<HighlightIcon sx={{ fontSize: theme.iconSize(1) }} />
							</IconButton>
						)}
						{(!playerStore.fullScreen && (
							<Tooltip title="Full screen">
								<IconButton onClick={playerStore.enterFullScreen}>
									<FullscreenIcon sx={{ fontSize: theme.iconSize(1) }} />
								</IconButton>
							</Tooltip>
						)) || (
							<Tooltip title="Exit full screen">
								<IconButton onClick={playerStore.exitFullScreen}>
									<FullscreenExitIcon sx={{ fontSize: theme.iconSize(1) }} />
								</IconButton>
							</Tooltip>
						)}
					</Box>
				</Stack>
			</Stack>
		</Fade>
	);
}

export default PlayerControls;
