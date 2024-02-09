import { useTheme } from '@emotion/react';
import { Box } from '@mui/material';
import React, { useEffect, useRef } from 'react';

function Thumbnail({ image, imageUrl, thumbnailUrl, title, crop }) {
	const thumbnailCanvasId = useRef(null);
	const theme = useTheme();

	useEffect(() => {
		if (thumbnailCanvasId.current == null) {
			return;
		}

		if (thumbnailUrl) {
			return;
		}

		thumbnailCanvasId.current.width = 80;
		thumbnailCanvasId.current.height = 80;

		var ctx = thumbnailCanvasId.current.getContext('2d');

		if (crop && crop.height != 0) {
			if (image) {
				ctx.drawImage(image, crop.x, crop.y, crop.width, crop.height, 0, 0, 80, 80);
			} else if (imageUrl) {
				let img = new Image();
				img.onload = function () {
					ctx.drawImage(img, crop.x, crop.y, crop.width, crop.height, 0, 0, 80, 80);
				};
				img.src = imageUrl;
			}
		} else if (title) {
			ctx.fillStyle = theme.palette.text.primary;
			ctx.font = theme.fontSize(0.9) + ' DejaVu';
			wrapText(ctx, title.toUpperCase(), thumbnailCanvasId.current.width, thumbnailCanvasId.current.height);
		}
	}, [crop, image, imageUrl, title]);

	const wrapText = (ctx, text, maxWidth, maxHeight) => {
		var words = text.split(' ');
		var line = '';
		let lineHeight = ctx.measureText('M').width + 5;
		let y = maxHeight / 2;

		for (var n = 0; n < words.length; n++) {
			var testLine = line + words[n] + ' ';
			var metrics = ctx.measureText(testLine);
			var textWidth = metrics.width;

			if (textWidth > maxWidth && n > 0) {
				ctx.fillText(line, (maxWidth - ctx.measureText(line).width) / 2, y);
				line = words[n] + ' ';
				y += lineHeight;
			} else {
				line = testLine;
			}
		}

		ctx.fillText(line, (maxWidth - ctx.measureText(line).width) / 2, y);
	};

	return (
		<Box
			sx={{
				border: theme.border(1, 'solid', theme.palette.text.primary),
				borderRadius: theme.spacing(0.5),
				width: theme.iconSize(3),
				height: theme.iconSize(3),
			}}
			component={thumbnailUrl ? 'img' : 'canvas'}
			src={thumbnailUrl}
			ref={thumbnailCanvasId}
		/>
	);
}

export default Thumbnail;
