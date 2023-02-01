import AddLink from '@mui/icons-material/AddLink';
import NoImageIcon from '@mui/icons-material/HideImage';
import ImageIcon from '@mui/icons-material/Image';
import OptionsIcon from '@mui/icons-material/Tune';
import { Box, SpeedDial, SpeedDialAction } from '@mui/material';
import React, { useRef, useState } from 'react';
import { useQueryClient } from 'react-query';
import Client from '../network/client';
import ReactQueryUtil from '../utils/react-query-util';
import TagAttachAnnotationMenu from './TagAttachAnnotationMenu';

function TagSpeedDial({ tag, hidden }) {
	const queryClient = useQueryClient();
	let [attachMenuAttributes, setAttachMenuAttributes] = useState(null);
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

	const onAttachAttributeClicked = (e) => {
		e.stopPropagation();
		setAttachMenuAttributes(
			attachMenuAttributes === null
				? {
						mouseX: e.clientX + 2,
						mouseY: e.clientY - 6,
				  }
				: null
		);
	};

	const menuClosed = (e) => {
		setAttachMenuAttributes(null);
	};

	return (
		<React.Fragment>
			{!hidden && attachMenuAttributes === null && (
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
					onKeyDown={(e) => e.stopPropagation()}
					onKeyUp={(e) => e.stopPropagation()}
				>
					<SpeedDialAction
						key="set-image"
						tooltipTitle="Set image"
						icon={<ImageIcon />}
						onKeyDown={(e) => e.stopPropagation()}
						onKeyUp={(e) => e.stopPropagation()}
						onClick={(e) => {
							changeTagImageClicked(e);
						}}
					/>
					<SpeedDialAction
						key="remove-image"
						tooltipTitle="Remove image"
						icon={<NoImageIcon />}
						onKeyDown={(e) => e.stopPropagation()}
						onKeyUp={(e) => e.stopPropagation()}
						onClick={(e) => {
							removeTagImageClicked(e);
						}}
					/>
					<SpeedDialAction
						key="manage-annotations"
						tooltipTitle="Manage annotations"
						icon={<AddLink />}
						onKeyDown={(e) => e.stopPropagation()}
						onKeyUp={(e) => e.stopPropagation()}
						onClick={(e) => {
							onAttachAttributeClicked(e);
						}}
					/>
				</SpeedDial>
			)}
			{!hidden && attachMenuAttributes !== null && (
				<TagAttachAnnotationMenu tag={tag} menu={attachMenuAttributes} onClose={menuClosed} />
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
