import ImageIcon from '@mui/icons-material/Image';
import OptionsIcon from '@mui/icons-material/Tune';
import { SpeedDial, SpeedDialAction, Typography } from '@mui/material';
import { useRef, useState } from 'react';
import Client from '../network/client';
import styles from './Tag.module.css';

function Tag({ tag, markActive, onTagSelected }) {
	let [optionsHidden, setOptionsHidden] = useState(true);
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

	const getTagClasses = () => {
		return styles.tag + ' ' + (tag.active && markActive ? styles.selected : styles.unselected);
	};

	const getCover = () => {
		if (tag.imageUrl) {
			return Client.buildStorageUrl(tag.imageUrl);
		} else {
			return 'empty';
		}
	};

	return (
		<div
			className={getTagClasses()}
			onClick={() => onTagSelected(tag)}
			onMouseEnter={() => setOptionsHidden(false)}
			onMouseLeave={() => setOptionsHidden(true)}
		>
			<img className={styles.image} src={getCover()} alt="" />
			<Typography className={styles.title} variant="h6" textAlign={'start'}>
				{tag.title}
			</Typography>
			<SpeedDial
				sx={{ '& .MuiFab-primary': { width: 40, height: 40, backgroundColor: 'rgba(0,0,0,0)' } }}
				hidden={optionsHidden}
				className={styles.tag_actions_button}
				ariaLabel="tag-actions"
				icon={<OptionsIcon />}
				onClick={(e) => e.stopPropagation()}
			>
				<SpeedDialAction
					key="set-image"
					tooltipTitle="Set tag image"
					icon={<ImageIcon />}
					onClick={(e) => {
						changeTagImageClicked(e);
					}}
				/>
			</SpeedDial>
			<input
				onClick={(e) => e.stopPropagation()}
				onChange={(e) => imageSelected(e)}
				accept="image/*"
				className={styles.choose_file_dialog}
				id="choose-file"
				type="file"
				ref={fileDialog}
			/>
		</div>
	);
}

export default Tag;
