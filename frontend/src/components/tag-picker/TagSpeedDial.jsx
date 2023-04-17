import AddLink from '@mui/icons-material/AddLink';
import { default as RemoveIcon } from '@mui/icons-material/Close';
import CopyIcon from '@mui/icons-material/ContentCopy';
import ImageIcon from '@mui/icons-material/Image';
import OptionsIcon from '@mui/icons-material/Tune';
import { SpeedDial, SpeedDialAction } from '@mui/material';
import React from 'react';

function TagSpeedDial({ tag, hidden, onManageAttributesClicked, onRemoveTagClicked, onManageImageClicked }) {
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
				</SpeedDial>
			)}
		</>
	);
}

export default TagSpeedDial;
