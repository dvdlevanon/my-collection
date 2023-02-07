import { Box } from '@mui/material';
import Client from '../network/client';

function Player({ item, isPreview }) {
	return (
		<Box
			component="video"
			height="100%"
			width="100%"
			playsInline
			muted
			autoPlay={isPreview}
			loop={isPreview}
			controls={!isPreview}
		>
			<source src={Client.buildFileUrl(isPreview ? item.previewUrl : item.url)} />
		</Box>
	);
}

export default Player;
