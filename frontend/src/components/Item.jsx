import { Tooltip } from '@mui/material';
import { useState } from 'react';
import { Link } from 'react-router-dom';
import Client from '../network/client';
import styles from './Item.module.css';
import ItemCoverIndicator from './ItemCoverIndicator';
import Player from './Player';

function Item({ item, preferPreview }) {
	let [mouseEnterMillis, setMouseEnterMillis] = useState(0);
	let [showCoverNavigator, setShowCoverNavigator] = useState(false);
	let [showPreview, setShowPreview] = useState(false);
	let [coverNumber, setCoverNumber] = useState(
		item.covers && item.covers.length > 0 ? Math.floor(item.covers.length / 2) : 0
	);

	const getCover = () => {
		if (item.covers && item.covers.length > 0) {
			return Client.buildFileUrl(item.covers[coverNumber].url);
		} else {
			return 'empty';
		}
	};

	const previewMode = () => {
		return preferPreview && item.previewUrl != '';
	};

	const mouseMove = (e) => {
		if (previewMode() || !item.covers) {
			return;
		}

		let nowMillis = Math.floor(Date.now());

		if (mouseEnterMillis == 0 || mouseEnterMillis > nowMillis - 250) {
			return;
		}

		let bounds = e.currentTarget.getBoundingClientRect();
		let x = e.clientX - bounds.left;
		setShowCoverNavigator(true);
		setCoverNumber(Math.floor(x / (bounds.width / item.covers.length)));
	};

	const mouseLeave = (e) => {
		if (previewMode()) {
			setShowPreview(false);
		} else {
			setMouseEnterMillis(0);
			setCoverNumber(item.covers && item.covers.length > 0 ? Math.floor(item.covers.length / 2) : 0);
			setShowCoverNavigator(false);
		}
	};

	const mouseEnter = (e) => {
		if (previewMode()) {
			setShowPreview(true);
		} else {
			setMouseEnterMillis(Math.floor(Date.now()));
		}
	};

	return (
		<Link
			className={styles.item}
			to={'item/' + item.id}
			onMouseLeave={(e) => mouseLeave(e)}
			onMouseMove={(e) => mouseMove(e)}
			onMouseEnter={(e) => mouseEnter(e)}
		>
			<img className={styles.image} src={getCover()} alt="" />
			{showCoverNavigator && item.covers && item.covers.length > 1 && (
				<div className={styles.cover_navigator}>
					{item.covers.map((cover, index) => {
						return (
							<ItemCoverIndicator
								key={cover.id}
								item={item}
								cover={cover}
								isHighlighted={coverNumber == index}
							/>
						);
					})}
				</div>
			)}
			{previewMode() && showPreview && <Player isPreview={true} item={item} />}
			<Tooltip title={item.title} arrow followCursor>
				<span className={styles.item_title}>{item.title}</span>
			</Tooltip>
		</Link>
	);
}

export default Item;
