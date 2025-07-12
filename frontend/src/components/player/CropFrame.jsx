import { useTheme } from '@emotion/react';
import { Box } from '@mui/material';
import { useEffect, useRef, useState } from 'react';
import CroppableImage from '../croppable-image/CroppableImage';

function CropFrame({ videoRef, isPlaying, width, height, onMouseMove, setCrop }) {
	const theme = useTheme();
	const thumbnailCanvasId = useRef(null);
	const [imageDataUrl, setImageDataUrl] = useState(null);

	useEffect(() => {
		if (thumbnailCanvasId.current == null) {
			return;
		}

		const canvas = thumbnailCanvasId.current;
		let ctx = canvas.getContext('2d');
		let animationFrameId;

		const drawVideoFrame = () => {
			if (!videoRef.current) {
				return;
			}

			const videoWidth = videoRef.current.videoWidth;
			const videoHeight = videoRef.current.videoHeight;

			if (canvas.width !== videoWidth || canvas.height !== videoHeight) {
				canvas.width = videoWidth;
				canvas.height = videoHeight;
			}

			if (videoRef.current && !videoRef.current.hasAttribute('crossorigin')) {
				videoRef.current.setAttribute('crossorigin', 'anonymous');
			}

			ctx.drawImage(videoRef.current, 0, 0, canvas.width, canvas.height);

			setImageDataUrl(canvas.toDataURL('image/png'));

			if (isPlaying) {
				animationFrameId = requestAnimationFrame(drawVideoFrame);
			}
		};

		drawVideoFrame();

		return () => {
			if (isPlaying) {
				cancelAnimationFrame(animationFrameId);
			}
		};
	}, [videoRef, isPlaying]);

	return (
		<Box
			borderRadius={theme.spacing(2)}
			height={height + 'px'}
			width={width + 'px'}
			position={'absolute'}
			zIndex={1}
			onMouseMove={onMouseMove}
			sx={{
				boxShadow: '3',
			}}
		>
			<Box
				position={'absolute'}
				component="canvas"
				ref={thumbnailCanvasId}
				top={0}
				left={0}
				height={height + 'px'}
				width={width + 'px'}
			></Box>
			{imageDataUrl && (
				<Box position={'absolute'} top={0} left={0} height={height + 'px'} width={width + 'px'}>
					<CroppableImage
						imageUrl={imageDataUrl}
						aspect={0}
						imageTitle={'test'}
						cropMode={true}
						showControls={false}
						onCropChange={setCrop}
						onImageLoaded={() => {}}
					/>
				</Box>
			)}
		</Box>
	);
}

export default CropFrame;
