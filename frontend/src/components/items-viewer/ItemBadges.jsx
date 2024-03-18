import { useTheme } from '@emotion/react';
import NoVideoIcon from '@mui/icons-material/VideocamOff';
import NoAudioIcon from '@mui/icons-material/VolumeOff';
import { Avatar, Box, IconButton, Tooltip } from '@mui/material';
import { Stack } from '@mui/system';
import React from 'react';
import CodecUtil from '../../utils/codec-utils';

function ItemBadges({ item }) {
	const theme = useTheme();

	return (
		<Box
			sx={{
				position: 'absolute',
				left: theme.spacing(1),
				top: theme.spacing(1),
				zIndex: 1000,
			}}
		>
			<Stack gap={theme.spacing(1)} flexDirection="row">
				{!CodecUtil.isVideoSupported(item.video_codec) && (
					<Tooltip title={'Video codec "' + item.video_codec + '" is not supported. click to convert'}>
						<IconButton
							onClick={(e) => {
								e.preventDefault();
							}}
						>
							<Avatar>
								<NoVideoIcon sx={{ fontSize: theme.iconSize(1) }} />
							</Avatar>
						</IconButton>
					</Tooltip>
				)}
				{!CodecUtil.isAudioSupported(item.audio_codec) && (
					<Tooltip title={'Audio codec "' + item.audio_codec + '" is not supported. click to convert'}>
						<IconButton
							onClick={(e) => {
								e.preventDefault();
							}}
						>
							<Avatar>
								<NoAudioIcon sx={{ fontSize: theme.iconSize(1) }} />
							</Avatar>
						</IconButton>
					</Tooltip>
				)}
			</Stack>
		</Box>
	);
}

export default ItemBadges;
