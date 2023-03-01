import NoVideoIcon from '@mui/icons-material/VideocamOff';
import NoAudioIcon from '@mui/icons-material/VolumeOff';
import { Avatar, Box, IconButton, Tooltip } from '@mui/material';
import { Stack } from '@mui/system';
import React from 'react';
import CodecUtil from '../../utils/codec-utils';

function ItemBadges({ item }) {
	return (
		<Box
			sx={{
				position: 'absolute',
				left: '10px',
				top: '10px',
				zIndex: 1000,
			}}
		>
			<Stack gap="10px" flexDirection="row">
				{!CodecUtil.isVideoSupported(item.videoCodec) && (
					<Tooltip title={'Video codec "' + item.videoCodec + '" is not supported. click to convert'}>
						<IconButton
							onClick={(e) => {
								e.preventDefault();
							}}
						>
							<NoVideoIcon color="bright" />
						</IconButton>
					</Tooltip>
				)}
				{!CodecUtil.isAudioSupported(item.audioCodec) && (
					<Tooltip title={'Audio codec "' + item.audioCodec + '" is not supported. click to convert'}>
						<IconButton
							onClick={(e) => {
								e.preventDefault();
							}}
						>
							<Avatar>
								<NoAudioIcon color="bright" />
							</Avatar>
						</IconButton>
					</Tooltip>
				)}
			</Stack>
		</Box>
	);
}

export default ItemBadges;
