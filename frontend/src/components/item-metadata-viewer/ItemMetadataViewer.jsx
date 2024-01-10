import { useTheme } from '@emotion/react';
import { Typography } from '@mui/material';
import React from 'react';
import BytesUtil from '../../utils/bytes-util';
import TimeUtil from '../../utils/time-utils';

function ItemMetadataViewer({ item }) {
	const theme = useTheme();

	return (
		<Typography variant="body2" color="bright.darker2" padding={theme.multiSpacing(0, 1)}>
			File Size: {BytesUtil.formatBytes(item.file_size, 2)}
			<br />
			Last Modified: {TimeUtil.formatEpochToDate(item.last_modified)}
			<br />
			Resolution: {item.width} * {item.height}
			<br />
			Video Codec: {item.video_codec}
			<br />
			Audio Codec: {item.audio_codec}
		</Typography>
	);
}

export default ItemMetadataViewer;
