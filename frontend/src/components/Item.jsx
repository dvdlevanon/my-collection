import { Tooltip } from '@mui/material';
import { useState } from 'react';
import { Link } from 'react-router-dom';
import Client from '../network/client';
import styles from './Item.module.css';
import ItemPreviewIndicator from './ItemPreviewIndicator';

function Item({ item }) {
	let [mouseEnterMillis, setMouseEnterMillis] = useState(0);
	let [showPreviewNavigator, setShowPreviewNavigator] = useState(false);
	let [previewNumber, setPreviewNumber] = useState(
		item.previews && item.previews.length > 0 ? Math.floor(item.previews.length / 2) : 0
	);

	const getCover = () => {
		if (item.previews && item.previews.length > 0) {
			return Client.buildStorageUrl(item.previews[previewNumber].url);
		} else {
			return 'empty';
		}
	};

	const mouseOver = (e) => {
		let nowMillis = Math.floor(Date.now());
		if (mouseEnterMillis == 0 || mouseEnterMillis > nowMillis - 250) {
			return;
		}

		if (!item.previews) {
			return;
		}

		let bounds = e.currentTarget.getBoundingClientRect();
		let x = e.clientX - bounds.left;
		setShowPreviewNavigator(true);
		setPreviewNumber(Math.floor(x / (bounds.width / item.previews.length)));
	};

	const mouseLeave = (e) => {
		setPreviewNumber(item.previews && item.previews.length > 0 ? Math.floor(item.previews.length / 2) : 0);
		setMouseEnterMillis(0);
		setShowPreviewNavigator(false);
	};

	const mouseEnter = (e) => {
		setMouseEnterMillis(Math.floor(Date.now()));
	};

	return (
		<Link
			className={styles.item}
			to={'item/' + item.id}
			onMouseLeave={(e) => mouseLeave(e)}
			onMouseMove={(e) => mouseOver(e)}
			onMouseEnter={(e) => mouseEnter(e)}
		>
			<img className={styles.image} src={getCover()} alt="" />
			{item.previews && showPreviewNavigator ? (
				<div className={styles.preview_navigator}>
					{item.previews.map((preview, index) => {
						return (
							<ItemPreviewIndicator
								key={preview.id}
								item={item}
								preview={preview}
								isHighlighted={previewNumber == index}
							/>
						);
					})}
				</div>
			) : (
				''
			)}
			<Tooltip title={item.title} arrow followCursor>
				<span className={styles.item_title}>{item.title}</span>
			</Tooltip>
		</Link>
	);
}

export default Item;
