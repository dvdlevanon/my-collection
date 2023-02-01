import AddLink from '@mui/icons-material/AddLink';
import { default as RemoveIcon } from '@mui/icons-material/Close';
import NoImageIcon from '@mui/icons-material/HideImage';
import ImageIcon from '@mui/icons-material/Image';
import OptionsIcon from '@mui/icons-material/Tune';
import { Box, SpeedDial, SpeedDialAction } from '@mui/material';
import React, { useRef } from 'react';
import { useQueryClient } from 'react-query';
import Client from '../network/client';
import ReactQueryUtil from '../utils/react-query-util';

function TagSpeedDial({ tag, hidden, onManageAttributesClicked, onRemoveTagClicked }) {
	const queryClient = useQueryClient();
	const fileDialog = useRef(null);

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
				queryClient.refetchQueries({ queryKey: ReactQueryUtil.TAGS_KEY });
			});
		});
	};

	const removeTagImageClicked = (e) => {
		Client.saveTag({ ...tag, imageUrl: 'none' }).then(() => {
			queryClient.refetchQueries({ queryKey: ReactQueryUtil.TAGS_KEY });
		});
	};

	return (
		<React.Fragment>
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
					onClick={(e) => e.stopPropagation()}
				>
					<SpeedDialAction
						key="set-image"
						tooltipTitle="Set image"
						icon={<ImageIcon />}
						onClick={(e) => {
							changeTagImageClicked(e);
						}}
					/>
					<SpeedDialAction
						key="remove-image"
						tooltipTitle="Remove image"
						icon={<NoImageIcon />}
						onClick={(e) => {
							removeTagImageClicked(e);
						}}
					/>
					<SpeedDialAction
						key="manage-annotations"
						tooltipTitle="Manage annotations"
						icon={<AddLink />}
						onClick={(e) => {
							onManageAttributesClicked(e);
						}}
					/>
					<SpeedDialAction
						key="remove-tag"
						tooltipTitle="Remove tag"
						icon={<RemoveIcon />}
						onClick={(e) => {
							e.stopPropagation();
							onRemoveTagClicked();
						}}
					/>
				</SpeedDial>
			)}
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
		</React.Fragment>
	);
}

export default TagSpeedDial;
