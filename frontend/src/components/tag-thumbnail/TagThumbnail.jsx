import GalleryIcon from '@mui/icons-material/Collections';
import DeleteIcon from '@mui/icons-material/Delete';
import ThumbnailIcon from '@mui/icons-material/Face3';
import { Box, Divider, ListItemIcon, Menu, MenuItem, Tooltip, Typography } from '@mui/material';
import { useState } from 'react';
import { Link, Link as RouterLink } from 'react-router-dom';
import seedrandom from 'seedrandom';
import Client from '../../utils/client';
import GalleryUrlParams from '../../utils/gallery-url-params';
import Thumbnail from '../thumbnail/Thumbnail';

function TagThumbnail({ tag, onTagClicked, onTagRemoved, onEditThumbnail, withRemoveOption }) {
	const [menuAchrosEl, setMenuAchrosEl] = useState(null);
	const [menuX, setMenuX] = useState(0);
	const [menuY, setMenuY] = useState(0);

	const onClick = (e) => {
		e.preventDefault();
		setMenuX(e.clientX);
		setMenuY(e.clientY);
		setMenuAchrosEl(e.target);
	};

	const closeMenu = () => {
		setMenuAchrosEl(null);
	};

	const getRandomImage = (images) => {
		let epochDay = Math.floor(Date.now() / 1000 / 60 / 60 / 24);
		let rand = seedrandom(epochDay + tag.id);
		let randomIndex = Math.floor(rand() * images.length);
		return images[randomIndex];
	};

	const getThumbnailCompnent = () => {
		if (!tag.images) {
			return <Thumbnail title={tag.title} />;
		}

		let imagesWithThumbnails = tag.images.filter((image) => {
			return image.thumbnail_rect && image.thumbnail_rect.height != 0;
		});

		if (imagesWithThumbnails.length == 0) {
			return <Thumbnail title={tag.title} />;
		}

		let image = getRandomImage(imagesWithThumbnails);
		return <Thumbnail crop={image.thumbnail_rect} imageUrl={Client.buildFileUrl(image.url)} />;
	};

	const getThumbnailCompnentWrapper = () => {
		return (
			<Tooltip title={tag.title}>
				<Box
					onContextMenu={onClick}
					onClick={() => {
						if (onTagClicked) {
							onTagClicked(tag);
						}
					}}
					sx={{
						cursor: 'pointer',
					}}
				>
					{getThumbnailCompnent()}
				</Box>
			</Tooltip>
		);
	};

	return (
		<>
			{onTagClicked == null && (
				<Link target="_blank" component={RouterLink} to={'/?' + GalleryUrlParams.buildUrlParams(tag.id)}>
					{getThumbnailCompnentWrapper()}
				</Link>
			)}
			{onTagClicked != null && getThumbnailCompnentWrapper()}
			<Menu
				open={menuAchrosEl != null}
				anchorEl={menuAchrosEl}
				onClose={() => setMenuAchrosEl(null)}
				anchorPosition={{ left: menuX, top: menuY }}
				anchorReference="anchorPosition"
			>
				<MenuItem disabled>
					<Typography variant="h5" color="white">
						{tag.title}
					</Typography>
				</MenuItem>
				<MenuItem
					onClick={() => {
						onEditThumbnail(tag);
						closeMenu();
					}}
				>
					<ListItemIcon>
						<ThumbnailIcon />
					</ListItemIcon>
					<Typography variant="h6" color="white">
						Set Thumbnail
					</Typography>
				</MenuItem>
				<Link target="_blank" component={RouterLink} to={'/?' + GalleryUrlParams.buildUrlParams(tag.id)}>
					<MenuItem onClick={closeMenu}>
						<ListItemIcon>
							<GalleryIcon />
						</ListItemIcon>
						<Typography variant="h6" color="white">
							Open in Gallery
						</Typography>
					</MenuItem>
				</Link>
				{withRemoveOption && <Divider />}
				{withRemoveOption && (
					<MenuItem
						onClick={() => {
							onTagRemoved(tag);
							closeMenu();
						}}
					>
						<ListItemIcon>
							<DeleteIcon />
						</ListItemIcon>
						<Typography variant="h6" color="white">
							Remove
						</Typography>
					</MenuItem>
				)}
			</Menu>
		</>
	);
}

export default TagThumbnail;
