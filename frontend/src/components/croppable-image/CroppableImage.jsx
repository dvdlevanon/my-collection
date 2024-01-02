import NoImageIcon from '@mui/icons-material/HideImage';
import { Box } from '@mui/material';
import React, { useEffect, useRef, useState } from 'react';
import ReactCrop from 'react-image-crop';

function CroppableImage({ imageUrl, imageTitle, cropMode, onCropChange, onImageLoaded }) {
	const [initialCrop, setInitialCrop] = useState(true);
	const [crop, setCrop] = useState(null);
	const [imageDimenssion, setImageDimenssion] = useState(null);
	const [originalImageDimenssion, setOriginalImageDimenssion] = useState(null);
	const imageHolderRef = useRef(null);
	const imageRef = useRef(null);

	useEffect(() => {
		const handleResize = (entries) => {
			for (let entry of entries) {
				const { width, height } = entry.contentRect;
				calculateImageSize(width, height);
			}
		};

		const observer = new ResizeObserver(handleResize);
		if (imageHolderRef.current) {
			observer.observe(imageHolderRef.current);
		}

		return () => {
			observer.disconnect();
		};
	}, [imageHolderRef.current]);

	const calculateImageSize = (width, height) => {
		if (!imageRef.current || imageRef.current.naturalHeight === 0) {
			return;
		}

		var aspectRatioImage = imageRef.current.naturalWidth / imageRef.current.naturalHeight;
		var aspectRatioContainer = width / height;

		var newWidth, newHeight;

		if (aspectRatioImage > aspectRatioContainer) {
			newWidth = width;
			newHeight = width / aspectRatioImage;
		} else {
			newHeight = height;
			newWidth = height * aspectRatioImage;
		}

		setImageDimenssion({ width: newWidth, height: newHeight });

		if (initialCrop) {
			setCrop({
				unit: 'px',
				x: newWidth / 2 - 100,
				y: newHeight / 2 - 100,
				height: 200,
				width: 200,
			});
		}
	};

	const getImageComponent = (e) => {
		return (
			<Box
				sx={{
					borderRadius: '5px',
					objectFit: 'contain',
					overflow: 'hidden',
					height: imageDimenssion ? imageDimenssion.height : 'auto',
					width: imageDimenssion ? imageDimenssion.width : 'auto',
				}}
				component="img"
				ref={imageRef}
				src={imageUrl}
				alt={imageTitle}
				loading="lazy"
				onLoad={() => {
					if (!imageRef.current) {
						return;
					}

					onImageLoaded(imageRef.current);
					setOriginalImageDimenssion({
						width: imageRef.current.naturalWidth,
						height: imageRef.current.naturalHeight,
					});

					if (!imageHolderRef.current) {
						return;
					}

					let imageBoundries = imageHolderRef.current.getBoundingClientRect();
					calculateImageSize(imageBoundries.width, imageBoundries.height);
				}}
			/>
		);
	};

	const cropChanged = (cropRect) => {
		setInitialCrop(false);
		setCrop(cropRect);

		if (imageRef.current == null) {
			return;
		}

		let imageBoundries = imageRef.current.getBoundingClientRect();
		let scaleX = originalImageDimenssion.width / imageBoundries.width;
		let scaleY = originalImageDimenssion.height / imageBoundries.height;
		var adjustedX = cropRect.x * scaleX;
		var adjustedY = cropRect.y * scaleY;
		var adjustedWidth = cropRect.width * scaleX;
		var adjustedHeight = cropRect.height * scaleY;

		onCropChange({
			x: adjustedX,
			y: adjustedY,
			width: adjustedWidth,
			height: adjustedHeight,
		});
	};

	return (
		<Box
			ref={imageHolderRef}
			sx={{
				display: 'flex',
				height: '100%',
				width: '100%',
				gap: '10px',
				justifyContent: 'center',
			}}
		>
			{(cropMode && (
				<ReactCrop aspect={1} crop={crop} onChange={cropChanged}>
					{getImageComponent()}
				</ReactCrop>
			)) ||
				getImageComponent()}
			{!imageUrl && (
				<Box
					sx={{
						position: 'absolute',
						'&:hover': {
							filter: 'brightness(120%)',
						},
						width: '100px',
						height: '100px',
						left: 0,
						right: 0,
						top: 0,
						bottom: 0,
						margin: 'auto',
						display: 'flex',
						flexDirection: 'column',
					}}
				>
					<NoImageIcon
						color="dark"
						sx={{
							fontSize: '100px',
						}}
					/>
				</Box>
			)}
		</Box>
	);
}

export default CroppableImage;
