import { useTheme } from '@emotion/react';
import CloseIcon from '@mui/icons-material/Close';
import { Box, Button, Dialog, DialogContent, DialogTitle, IconButton, Stack, Typography } from '@mui/material';
import React, { useEffect, useRef, useState } from 'react';
import 'react-image-crop/dist/ReactCrop.css';
import { useQuery, useQueryClient } from 'react-query';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import TagsUtil from '../../utils/tags-util';
import CroppableImage from '../croppable-image/CroppableImage';
import TagImageTypeSelector from '../tag-picker/TagImageTypeSelector';
import Thumbnail from '../thumbnail/Thumbnail';

function ManageTagImageDialog({ tag, autoThumbnailMode, onClose }) {
	const queryClient = useQueryClient();
	const [tit, setTit] = useState(null);
	const [updatedTag, setUpdatedTag] = useState(tag);
	const [thumbnailMode, setThumbnailMode] = useState(autoThumbnailMode);
	const [thumbnailCrop, setThumbnailCrop] = useState(null);
	const [image, setImage] = useState(null);
	const fileDialog = useRef(null);
	const theme = useTheme();
	const titsQuery = useQuery({
		queryKey: ReactQueryUtil.TAG_IMAGE_TYPES_KEY,
		queryFn: Client.getTagImageTypes,
		onSuccess: (tits) => {
			let lastTit = localStorage.getItem('manage_tag_image_last_tit');

			if (lastTit) {
				setTit(tits.find((cur) => cur.id == lastTit));
			} else if (!tit) {
				setTit(tits[0]);
			}
		},
	});

	useEffect(() => {
		let lastTit = localStorage.getItem('manage_tag_image_last_tit');
		if (lastTit && titsQuery.isSuccess) {
			setTit(titsQuery.data.find((cur) => cur.id == lastTit));
		}
	});

	const imageFromClipboardClicked = async (e) => {
		e.stopPropagation();
		const imageUrl = await navigator.clipboard.readText();
		if (!imageUrl.startsWith('http')) {
			alert('Invalid clipboard data ' + imageUrl);
			return;
		}

		if (TagsUtil.hasTagImage(tag, tit)) {
			alert('Remove current image first');
			return;
		}

		let fileName = TagsUtil.tagTitleToFileName(tag.title) + '.' + imageUrl.split('.').pop();

		Client.uploadFileFromUrl(`tags-image-types/${tag.id}/${tit.id}/${fileName}`, imageUrl, (fileUrl) => {
			updateTagImage(tag, tit, fileUrl.url);
		});
	};

	const changeTagImageClicked = (e) => {
		e.stopPropagation();
		fileDialog.current.value = '';
		fileDialog.current.click();
	};

	const updateTagImage = (tag, tit, fileUrl) => {
		if (tag.images && tag.images.some((image) => image.imageType === tit.id)) {
			return;
		}

		let updatedTag = { ...tag };
		updatedTag.images = updatedTag.images || [];
		updatedTag.images.push({
			url: fileUrl,
			image_nonce: Date.now(),
			tag_id: tag.id,
			imageType: tit.id,
		});

		Client.saveTag(updatedTag).then(() => {
			setUpdatedTag(updatedTag);
			queryClient.refetchQueries({ queryKey: ReactQueryUtil.TAGS_KEY });
		});
	};

	const updateImageThumbnail = (tag, tit, thumbnail) => {
		let updatedTag = { ...tag };
		let image = updatedTag.images.find((image) => image.imageType === tit.id);
		image.thumbnail_rect = {
			x: Math.floor(thumbnail.x),
			y: Math.floor(thumbnail.y),
			width: Math.floor(thumbnail.width),
			height: Math.floor(thumbnail.height),
		};

		Client.updateTagImage(tag.id, image).then(() => {
			setUpdatedTag(updatedTag);
			queryClient.refetchQueries({ queryKey: ReactQueryUtil.TAGS_KEY });
		});
	};

	const imageSelected = (e) => {
		if (fileDialog.current.files.length !== 1) {
			return;
		}

		Client.uploadFile(`tags-image-types/${tag.id}/${tit.id}`, fileDialog.current.files[0], (fileUrl) => {
			updateTagImage(updatedTag, tit, fileUrl.url);
		});
	};

	const removeTagImageClicked = (e) => {
		e.stopPropagation();

		if (!TagsUtil.hasTagImage(tag, tit)) {
			alert('No image to remove');
			return;
		}

		Client.removeTagImageFromTag(tag.id, tit.id).then(() => {
			let updatedTag = { ...tag };
			updatedTag.images = updatedTag.images || [];
			updatedTag.images = updatedTag.images.filter((image) => image.imageType !== tit.id);
			setUpdatedTag(updatedTag);
			queryClient.refetchQueries({ queryKey: ReactQueryUtil.TAGS_KEY });
		});
	};

	const getRegularButtons = (e) => {
		return (
			<>
				<Button
					variant="outlined"
					onClick={(e) => {
						updateImageThumbnail(tag, tit, { x: 0, y: 0, height: 0, width: 0 });
					}}
				>
					Remove Thumbnail
				</Button>
				<Button
					disabled={!TagsUtil.hasImage(updatedTag)}
					variant="outlined"
					onClick={(e) => {
						removeTagImageClicked(e);
					}}
				>
					Remove Image
				</Button>
				<Button
					onClick={(e) => {
						changeTagImageClicked(e);
					}}
					variant="outlined"
				>
					Upload Image
				</Button>
				<Button
					onClick={(e) => {
						imageFromClipboardClicked(e);
					}}
					variant="outlined"
				>
					Image From Clipboard
				</Button>
				<Button
					variant="outlined"
					onClick={(e) => {
						setThumbnailMode(true);
					}}
				>
					Set Thumbnail
				</Button>
			</>
		);
	};
	return (
		<Dialog
			onClose={(e, reason) => {
				e.stopPropagation();
				if (reason === 'backdropClick' || reason === 'escapeKeyDown') {
					onClose();
				}
			}}
			open={true}
			fullWidth={true}
			maxWidth={'xl'}
			PaperProps={{ sx: { maxHeight: '80vh', minHeight: '80vh', height: '80vh' } }}
		>
			<DialogTitle onClick={(e) => e.stopPropagation()}>
				<Stack flexDirection="row" gap={theme.spacing(2)}>
					<Typography variant="h6">Set Image for {updatedTag.title}</Typography>
					{titsQuery.isSuccess && (
						<TagImageTypeSelector
							disabled={thumbnailMode}
							tits={titsQuery.data}
							tit={tit}
							onTitChanged={(tit) => {
								localStorage.setItem('manage_tag_image_last_tit', tit.id);
								setTit(tit);
							}}
						/>
					)}
				</Stack>
				<IconButton
					sx={{
						position: 'absolute',
						top: '0',
						right: '0',
						margin: theme.spacing(1),
					}}
					onClick={(e) => {
						e.stopPropagation();
						onClose();
					}}
				>
					<CloseIcon sx={{ fontSize: theme.iconSize(1) }} />
				</IconButton>
			</DialogTitle>
			<DialogContent
				sx={{
					display: 'flex',
					flexDirection: 'column',
					gap: theme.spacing(0.5),
					padding: theme.spacing(0.5),
				}}
				onClick={(e) => e.stopPropagation()}
			>
				<Box
					sx={{
						width: '100%',
						height: '100%',
						overflow: 'hidden',
					}}
				>
					<Stack
						flexDirection="row"
						gap={theme.spacing(1)}
						sx={{
							width: '100%',
							height: '100%',
							justifyContent: 'center',
						}}
					>
						<CroppableImage
							imageUrl={TagsUtil.getTagImageUrl(updatedTag, tit, true)}
							imageTitle={updatedTag.title}
							cropMode={thumbnailMode}
							onCropChange={setThumbnailCrop}
							onImageLoaded={setImage}
							onCropCanceled={() => {
								setThumbnailMode(false);
								setThumbnailCrop(null);
							}}
							onCropDone={() => {
								updateImageThumbnail(tag, tit, thumbnailCrop);
								setThumbnailMode(false);
								setThumbnailCrop(null);
							}}
						/>
						<Thumbnail
							imageUrl={TagsUtil.getTagImageUrl(updatedTag, tit, true)}
							image={thumbnailMode ? image : null}
							crop={thumbnailMode ? thumbnailCrop : TagsUtil.getTagImageThumbnailRect(updatedTag, tit)}
							title={thumbnailMode ? '' : 'Thum- bnail'}
						/>
					</Stack>
				</Box>
				<Box
					sx={{
						display: 'flex',
						flexDirection: 'row',
						justifyContent: 'center',
						gap: theme.spacing(1),
						padding: theme.spacing(1),
					}}
					onClick={(e) => e.stopPropagation()}
				>
					{!thumbnailMode && getRegularButtons()}
				</Box>
				<Box
					component="input"
					onClick={(e) => e.stopPropagation()}
					onChange={(e) => imageSelected(e)}
					accept="image/*"
					id="choose-file"
					type="file"
					sx={{
						display: 'none',
					}}
					ref={fileDialog}
				/>
			</DialogContent>
		</Dialog>
	);
}

export default ManageTagImageDialog;
