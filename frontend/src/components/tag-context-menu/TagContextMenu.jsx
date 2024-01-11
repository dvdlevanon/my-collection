import { useTheme } from '@emotion/react';
import AddLink from '@mui/icons-material/AddLink';
import CopyIcon from '@mui/icons-material/ContentCopy';
import { default as RemoveIcon } from '@mui/icons-material/Delete';
import ThumbnailIcon from '@mui/icons-material/Face3';
import ImageIcon from '@mui/icons-material/Image';
import GalleryIcon from '@mui/icons-material/OpenInNew';
import { Divider, ListItemIcon, Menu, MenuItem, Typography } from '@mui/material';
import React from 'react';
import { useQuery } from 'react-query';
import { Link, Link as RouterLink } from 'react-router-dom';
import Client from '../../utils/client';
import GalleryUrlParams from '../../utils/gallery-url-params';
import ReactQueryUtil from '../../utils/react-query-util';

function TagContextMenu({
	tag,
	menuAnchorEl,
	onClose,
	menuPosition,
	onManageImageClicked,
	onManageAttributesClicked,
	onRemoveTagClicked,
	onEditThumbnail,
	withRemoveOption,
	withManageAttributesClicked,
}) {
	const theme = useTheme();
	const tagCustomCommandsQuery = useQuery(ReactQueryUtil.tagCustomCommands(tag.parentId), () =>
		Client.getTagCustomCommands(tag.parentId)
	);

	const closeEndCall = (e, callIt) => {
		onClose();
		callIt(e);
	};

	const handleCustomCommand = (e, command) => {
		let commands = command.type.split(',');

		for (let i = 0; i < commands.length; i++) {
			if (commands[i] == 'search web') {
				let url = command.arg.replace('${tag_title}', tag.title).replace(' ', '+');
				window.open(url, '_newtab');
			} else if (commands[i] == 'open-manage-tag-image-dialog') {
				onManageImageClicked(e);
			} else {
				console.log('Unknown command type ' + command.type);
			}
		}
	};

	return (
		<Menu
			open={true}
			anchorEl={menuAnchorEl}
			onClose={onClose}
			anchorPosition={menuPosition}
			anchorReference="anchorPosition"
		>
			<MenuItem disabled>
				<Typography variant="h5" color="white">
					{tag.title}
				</Typography>
			</MenuItem>
			<MenuItem
				onClick={() => {
					navigator.clipboard.writeText(tag.title);
					onClose();
				}}
			>
				<ListItemIcon>
					<CopyIcon sx={{ fontSize: theme.iconSize(1) }} />
				</ListItemIcon>
				Copy Title
			</MenuItem>
			<Link target="_blank" component={RouterLink} to={'/?' + GalleryUrlParams.buildUrlParams(tag.id)}>
				<MenuItem onClick={onClose} sx={{ color: 'white' }}>
					<ListItemIcon>
						<GalleryIcon sx={{ fontSize: theme.iconSize(1) }} />
					</ListItemIcon>
					Open in Gallery
				</MenuItem>
			</Link>
			<MenuItem onClick={(e) => closeEndCall(e, onManageImageClicked)}>
				<ListItemIcon>
					<ImageIcon sx={{ fontSize: theme.iconSize(1) }} />
				</ListItemIcon>
				Image Options...
			</MenuItem>
			<MenuItem
				onClick={(e) => {
					closeEndCall(e, onEditThumbnail);
				}}
			>
				<ListItemIcon>
					<ThumbnailIcon sx={{ fontSize: theme.iconSize(1) }} />
				</ListItemIcon>
				Set Thumbnail
			</MenuItem>
			{withManageAttributesClicked && (
				<MenuItem onClick={(e) => closeEndCall(e, onManageAttributesClicked)}>
					<ListItemIcon>
						<AddLink sx={{ fontSize: theme.iconSize(1) }} />
					</ListItemIcon>
					Manage annotations
				</MenuItem>
			)}
			{withRemoveOption && (
				<MenuItem onClick={(e) => closeEndCall(e, onRemoveTagClicked)}>
					<ListItemIcon>
						<RemoveIcon sx={{ fontSize: theme.iconSize(1) }} />
					</ListItemIcon>
					Remove
				</MenuItem>
			)}
			{tagCustomCommandsQuery.isSuccess && tagCustomCommandsQuery.data.length > 0 && <Divider />}
			{tagCustomCommandsQuery.isSuccess &&
				tagCustomCommandsQuery.data.length > 0 &&
				tagCustomCommandsQuery.data.map((command) => {
					return (
						<MenuItem
							key={command.id}
							onClick={(e) => {
								handleCustomCommand(e, command);
								onClose();
							}}
						>
							<ListItemIcon>
								{
									<img
										src={command.icon}
										alt="icon"
										style={{ width: theme.iconSize(1), height: theme.iconSize(1) }}
									/>
								}
							</ListItemIcon>
							{command.tooltip}
						</MenuItem>
					);
				})}
		</Menu>
	);
}

export default TagContextMenu;
