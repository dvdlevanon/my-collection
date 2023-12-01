import CloseIcon from '@mui/icons-material/Close';
import NoImageIcon from '@mui/icons-material/HideImage';
import { Box, Button, Dialog, DialogContent, DialogTitle, IconButton, Stack, Typography } from '@mui/material';
import React, { useEffect, useRef, useState } from 'react';
import { useQuery, useQueryClient } from 'react-query';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import TagsUtil from '../../utils/tags-util';
import TagImageTypeSelector from '../tag-picker/TagImageTypeSelector';

function ManageTagImageDialog({ tag, onClose }) {
	const queryClient = useQueryClient();
	const [tit, setTit] = useState(null);
	const [updatedTag, setUpdatedTag] = useState(tag);
	const titsQuery = useQuery({
		queryKey: ReactQueryUtil.TAG_IMAGE_TYPES_KEY,
		queryFn: Client.getTagImageTypes,
		onSuccess: (tits) => {
			let lastTit = localStorage.getItem('manage_tag_image_last_tit');

			if (lastTit) {
				setTit(titsQuery.data.find((cur) => cur.id == lastTit));
			} else if (!tit) {
				setTit(tits[0]);
			}
		},
	});

	const fileDialog = useRef(null);

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
		updatedTag.images.push({ url: fileUrl, tag_id: tag.id, imageType: tit.id });

		Client.saveTag(updatedTag).then(() => {
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
				<Stack flexDirection="row" gap="20px">
					<Typography variant="h6">Set Image for {updatedTag.title}</Typography>
					{titsQuery.isSuccess && (
						<TagImageTypeSelector
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
						top: '0px',
						right: '0px',
						margin: '10px',
					}}
					onClick={(e) => {
						e.stopPropagation();
						onClose();
					}}
				>
					<CloseIcon />
				</IconButton>
			</DialogTitle>
			<DialogContent
				sx={{
					display: 'flex',
					flexDirection: 'column',
					gap: '5px',
					padding: '5px',
				}}
				onClick={(e) => e.stopPropagation()}
			>
				<Box
					sx={{
						position: 'relative',
						display: 'flex',
						justifyContent: 'center',
						objectFit: 'contain',
						overflow: 'hidden',
						borderRadius: '5px',
						height: '100%',
					}}
				>
					<Box></Box>
					<Box
						sx={{
							borderRadius: '5px',
							objectFit: 'contain',
							overflow: 'hidden',
						}}
						component="img"
						src={TagsUtil.getTagImageUrl(updatedTag, tit, true)}
						alt={updatedTag.title}
						loading="lazy"
					/>
					{!TagsUtil.hasImage(updatedTag) && (
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
				<Box
					sx={{
						display: 'flex',
						flexDirection: 'row',
						justifyContent: 'center',
						gap: '10px',
						padding: '10px',
					}}
					onClick={(e) => e.stopPropagation()}
				>
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
