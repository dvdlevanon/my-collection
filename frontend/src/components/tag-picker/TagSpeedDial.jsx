import AddLink from '@mui/icons-material/AddLink';
import CopyIcon from '@mui/icons-material/ContentCopy';
import { default as RemoveIcon } from '@mui/icons-material/Delete';
import ImageIcon from '@mui/icons-material/Image';
import OptionsIcon from '@mui/icons-material/Tune';
import { SpeedDial, SpeedDialAction } from '@mui/material';
import React from 'react';
import { useQuery } from 'react-query';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';

function TagSpeedDial({ tag, hidden, onManageAttributesClicked, onRemoveTagClicked, onManageImageClicked }) {
	const tagCustomCommandsQuery = useQuery(ReactQueryUtil.tagCustomCommands(tag.parentId), () =>
		Client.getTagCustomCommands(tag.parentId)
	);

	const handleCustomCommand = (e, command) => {
		e.preventDefault();
		e.stopPropagation();

		if (command.type == 'search web') {
			let url = command.arg.replace('${tag_title}', tag.title).replace(' ', '+');
			window.open(url, '_newtab');
			onManageImageClicked();
		} else {
			console.log('Unknown command type ' + command.type);
		}
	};

	return (
		<>
			{!hidden && (
				<SpeedDial
					sx={{
						position: 'absolute',
						bottom: '0px',
						right: '0px',
						padding: '5px',
						'& .MuiFab-primary': {
							width: 40,
							height: 40,
							backgroundColor: 'primary.main',
						},
					}}
					ariaLabel="tag-actions"
					icon={<OptionsIcon />}
					onClick={(e) => {
						e.preventDefault();
						e.stopPropagation();
					}}
				>
					<SpeedDialAction
						key="copy-name"
						tooltipTitle="Copy title to clipboard"
						icon={<CopyIcon />}
						onClick={(e) => {
							navigator.clipboard.writeText(tag.title);
						}}
					/>

					<SpeedDialAction
						key="manage-image"
						tooltipTitle="Image options"
						icon={<ImageIcon />}
						onClick={(e) => {
							onManageImageClicked();
						}}
					/>

					<SpeedDialAction
						key="manage-annotations"
						tooltipTitle="Manage annotations"
						icon={<AddLink />}
						onClick={(e) => {
							e.preventDefault();
							e.stopPropagation();
							onManageAttributesClicked(e);
						}}
					/>
					<SpeedDialAction
						key="remove-tag"
						tooltipTitle="Remove tag"
						icon={<RemoveIcon />}
						onClick={(e) => {
							e.preventDefault();
							e.stopPropagation();
							onRemoveTagClicked();
						}}
					/>
					{tagCustomCommandsQuery.isSuccess &&
						tagCustomCommandsQuery.data.map((command) => {
							return (
								<SpeedDialAction
									key={command.id}
									tooltipTitle={command.tooltip}
									icon={
										<img src={command.icon} alt="icon" style={{ width: '24px', height: '24px' }} />
									}
									onClick={(e) => handleCustomCommand(e, command)}
								/>
							);
						})}
				</SpeedDial>
			)}
		</>
	);
}

export default TagSpeedDial;
