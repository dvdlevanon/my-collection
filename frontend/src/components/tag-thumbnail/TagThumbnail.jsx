import GalleryIcon from '@mui/icons-material/Collections';
import DeleteIcon from '@mui/icons-material/Delete';
import ThumbnailIcon from '@mui/icons-material/Face3';
import { Box, Divider, ListItemIcon, Menu, MenuItem, Tooltip, Typography } from '@mui/material';
import { useState } from 'react';
import { Link, Link as RouterLink } from 'react-router-dom';
import GalleryUrlParams from '../../utils/gallery-url-params';
import TagsUtil from '../../utils/tags-util';

function TagThumbnail({ tag, onTagRemoved, onEditThumbnail }) {
	const [menuAchrosEl, setMenuAchrosEl] = useState(null);
	const [menuX, setMenuX] = useState(0);
	const [menuY, setMenuY] = useState(0);

	const onClick = (e) => {
		setMenuX(e.clientX);
		setMenuY(e.clientY);
		setMenuAchrosEl(e.target);
	};

	const closeMenu = () => {
		setMenuAchrosEl(null);
	};

	return (
		<>
			<Tooltip title={tag.title}>
				<Box
					width="70px"
					height="70px"
					component="img"
					src={TagsUtil.getTagImageUrl(tag, null, false)}
					alt={tag.title}
					loading="lazy"
					onClick={onClick}
					sx={{
						cursor: 'pointer',
					}}
				></Box>
			</Tooltip>
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
				<MenuItem onClick={closeMenu}>
					<ListItemIcon>
						<GalleryIcon />
					</ListItemIcon>
					<Link target="_blank" component={RouterLink} to={'/?' + GalleryUrlParams.buildUrlParams(tag.id)}>
						<Typography variant="h6" color="white">
							Open in Gallery
						</Typography>
					</Link>
				</MenuItem>
				<Divider />
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
			</Menu>
		</>
	);
}

export default TagThumbnail;
