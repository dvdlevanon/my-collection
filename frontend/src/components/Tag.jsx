import ImageIcon from '@mui/icons-material/Image';
import OptionsIcon from '@mui/icons-material/Tune';
import { Box, SpeedDial, SpeedDialAction, Typography } from '@mui/material';
import { useRef, useState } from 'react';
import Client from '../network/client';
import styles from './Tag.module.css';

function Tag({ tag, size, onTagSelected }) {
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
		if (size != 'small') {
			return styles.tag + ' ' + styles.big_tag;
		} else {
			return styles.tag;
		}
	};

	const getImageUrl = () => {
		if (tag.imageUrl) {
			return Client.buildFileUrl(tag.imageUrl);
		} else {
			return 'empty';
		}
	};

	return (
		<Box
			className={getTagClasses()}
			onClick={() => onTagSelected(tag)}
			onMouseEnter={() => setOptionsHidden(false)}
			onMouseLeave={() => setOptionsHidden(true)}
		>
			{size != 'small' && <img className={styles.image} src={getImageUrl()} alt={tag.title} loading="lazy" />}
			<Typography
				className={styles.title}
				sx={{
					'&:hover': {
						textDecoration: 'underline',
					},
				}}
				variant="caption"
				textAlign={'start'}
			>
				{tag.title}
			</Typography>
			{size != 'small' && (
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
		</Box>
	);
}

export default Tag;
