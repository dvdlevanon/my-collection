import { Box } from '@mui/material';
import React, { useEffect, useRef } from 'react';

function Thumbnail({ image, crop }) {
	const thumbnailCanvasId = useRef(null);

	useEffect(() => {
		if (thumbnailCanvasId.current == null) {
			return;
		}

		if (crop == null) {
			return;
		}

		thumbnailCanvasId.current.width = 80;
		thumbnailCanvasId.current.height = 80;

		var ctx = thumbnailCanvasId.current.getContext('2d');

		ctx.drawImage(image, crop.x, crop.y, crop.width, crop.height, 0, 0, 80, 80);
	}, [thumbnailCanvasId.current, crop]);

	return (
		<Box
			sx={{
				border: 'white 1px solid',
				borderRadius: '5px',
				width: '80px',
				height: '80px',
			}}
			component="canvas"
			ref={thumbnailCanvasId}
		/>
	);
}

export default Thumbnail;
