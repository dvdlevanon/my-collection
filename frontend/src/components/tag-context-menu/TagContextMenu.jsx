import AddLink from '@mui/icons-material/AddLink';
import CopyIcon from '@mui/icons-material/ContentCopy';
import { default as RemoveIcon } from '@mui/icons-material/Delete';
import ImageIcon from '@mui/icons-material/Image';
import { Divider, ListItemIcon, Menu, MenuItem, Typography } from '@mui/material';
import React from 'react';
import { useQuery } from 'react-query';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';

function TagContextMenu({
	tag,
	menuAnchorEl,
	onClose,
	menuPosition,
	onManageImageClicked,
	onManageAttributesClicked,
	onRemoveTagClicked,
}) {
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
					<CopyIcon />
				</ListItemIcon>
				Copy Title
			</MenuItem>
			<MenuItem onClick={(e) => closeEndCall(e, onManageImageClicked)}>
				<ListItemIcon>
					<ImageIcon />
				</ListItemIcon>
				Image Options...
			</MenuItem>
			<MenuItem onClick={(e) => closeEndCall(e, onManageAttributesClicked)}>
				<ListItemIcon>
					<AddLink />
				</ListItemIcon>
				Manage annotations
			</MenuItem>
			<MenuItem onClick={(e) => closeEndCall(e, onRemoveTagClicked)}>
				<ListItemIcon>
					<RemoveIcon />
				</ListItemIcon>
				Remove
			</MenuItem>
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
								{<img src={command.icon} alt="icon" style={{ width: '24px', height: '24px' }} />}
							</ListItemIcon>
							{command.tooltip}
						</MenuItem>
					);
				})}
		</Menu>
	);
}

export default TagContextMenu;
