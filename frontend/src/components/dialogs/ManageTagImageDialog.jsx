import CloseIcon from '@mui/icons-material/Close';
import NoImageIcon from '@mui/icons-material/HideImage';
import { Box, Button, Dialog, DialogContent, DialogTitle, IconButton } from '@mui/material';
import React, { useRef } from 'react';
import { QueryClient, useQueryClient } from 'react-query';
import Client from '../../network/client';
import ReactQueryUtil from '../../utils/react-query-util';

function ManageTagImageDialog({ tag, onClose }) {
	const queryClient = useQueryClient();
	const fileDialog = useRef(null);

	const getImageUrl = () => {
		if (hasImage()) {
			return Client.buildFileUrl(tag.imageUrl);
		} else {
			return Client.buildFileUrl(Client.buildInternalStoragePath('tags-image/none/1.jpg'));
		}
	};

	const hasImage = () => {
		return tag.imageUrl && tag.imageUrl != 'none';
	};

	const changeTagImageClicked = (e) => {
		e.stopPropagation();
		fileDialog.current.value = '';
		fileDialog.current.click();
	};

	const imageSelected = (e) => {
		if (fileDialog.current.files.length != 1) {
			return;
		}

		Client.uploadFile(`tags-image/${tag.id}`, fileDialog.current.files[0], (fileUrl) => {
			Client.saveTag({ ...tag, imageUrl: fileUrl.url }).then(() => {
				onClose();
				QueryClient.refetchQueries({ queryKey: ReactQueryUtil.TAGS_KEY });
			});
		});
	};

	const removeTagImageClicked = (e) => {
		e.stopPropagation();
		Client.saveTag({ ...tag, imageUrl: 'none' }).then(() => {
			queryClient.refetchQueries({ queryKey: ReactQueryUtil.TAGS_KEY });
		});
	};

	return (
		<Dialog
			onClose={(e, reason) => {
				e.stopPropagation();
				if (reason == 'backdropClick' || reason == 'escapeKeyDown') {
					onClose();
				}
			}}
			open={true}
			fullWidth={true}
			maxWidth={'xl'}
		>
			<DialogTitle onClick={(e) => e.stopPropagation()}>
				Set Image for {tag.title}
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
					}}
				>
					<Box
						sx={{
							borderRadius: '5px',
							objectFit: 'contain',
							overflow: 'hidden',
						}}
						component="img"
						src={getImageUrl()}
						alt={tag.title}
						loading="lazy"
					/>
					{!hasImage() && (
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
						disabled={!hasImage()}
						variant="outlined"
						onClick={(e) => {
							if (!hasImage()) {
								return;
							}
							removeTagImageClicked(e);
							onClose();
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
						Set Image
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
