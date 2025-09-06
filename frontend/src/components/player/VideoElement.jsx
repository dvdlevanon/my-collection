import { useTheme } from '@emotion/react';
import { Box } from '@mui/material';
import Client from '../../utils/client';
import { usePlayerActionStore } from './PlayerActionStore';
import { usePlayerStore } from './PlayerStore';
import { useSubtitleStore } from './subtitles/SubtitlesStore';

function VideoElement({ videoController }) {
	const theme = useTheme();
	const playerStore = usePlayerStore();
	const playerActionStore = usePlayerActionStore();
	const subtitleStore = useSubtitleStore();

	return (
		playerStore.url && (
			<Box
				borderRadius={theme.spacing(2)}
				sx={{
					boxShadow: '3',
				}}
				component="video"
				crossOrigin="anonymous"
				height="100%"
				width="100%"
				controls={false}
				playsInline
				autoPlay={true}
				loop={false}
				ref={videoController.videoElement}
				onClick={() => {
					playerStore.togglePlay();
					playerActionStore.closeAll();
				}}
				onEnded={playerStore.videoFinished}
				onDoubleClick={playerStore.toggleFullScreen}
				onTimeUpdate={(e) => {
					playerStore.videoTimeUpdate(e.target.currentTime);
				}}
				onLoadedMetadata={(e) => {
					playerStore.videoLoadedMetadata(e.target.duration);
				}}
				onMouseMove={() => playerStore.showControls(true)}
			>
				<source src={Client.buildFileUrl(playerStore.url)} />
			</Box>
		)
	);
}

export default VideoElement;
