import AddLink from '@mui/icons-material/AddLink';
import ImageIcon from '@mui/icons-material/Image';
import RemoveLink from '@mui/icons-material/LinkOff';
import OptionsIcon from '@mui/icons-material/Tune';
import { SpeedDial, SpeedDialAction } from '@mui/material';
import React, { useRef, useState } from 'react';
import Client from '../network/client';
import TagAttachAnnotationMenu from './TagAttachAnnotationMenu';
import styles from './TagSpeedDial.module.css';

function TagSpeedDial({ tag }) {
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
			console.log(fileUrl);
			tag.imageUrl = fileUrl.url;
			Client.saveTag(tag);
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

	const onDettachAttributeClicked = (e) => {};

	const menuClosed = (e) => {
		setAttachMenuAttributes(null);
	};

	return (
		<React.Fragment>
			<SpeedDial
				sx={{ '& .MuiFab-primary': { width: 40, height: 40, backgroundColor: 'rgba(0,0,0,0)' } }}
				className={styles.tag_actions_button}
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
					key="attach-attribute"
					tooltipTitle="Attach attribute"
					icon={<AddLink />}
					onClick={(e) => {
						onAttachAttributeClicked(e);
					}}
				/>
				<SpeedDialAction
					key="remove-attribute"
					tooltipTitle="Remove attribute"
					icon={<RemoveLink />}
					onClick={(e) => {
						onDettachAttributeClicked(e);
					}}
				/>
			</SpeedDial>
			{attachMenuAttributes !== null && (
				<TagAttachAnnotationMenu tag={tag} menu={attachMenuAttributes} onClose={menuClosed} />
			)}
			<input
				onClick={(e) => e.stopPropagation()}
				onChange={(e) => imageSelected(e)}
				accept="image/*"
				className={styles.choose_file_dialog}
				id="choose-file"
				type="file"
				ref={fileDialog}
			/>
		</React.Fragment>
	);
}

export default TagSpeedDial;
