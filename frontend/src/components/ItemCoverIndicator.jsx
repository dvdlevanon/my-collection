import React, { useState } from 'react';
import styles from './ItemCoverIndicator.module.css';

function ItemCoverIndicator({ item, cover, isHighlighted }) {
	let [optionsHidden, setOptionsHidden] = useState(true);

	const clicked = (e) => {
		e.stopPropagation();
		console.log('cliecked');
	};

	return (
		<span
			className={styles.cover_indicator + ' ' + (isHighlighted ? styles.selected : styles.unselected)}
			key={cover.id}
			onMouseEnter={() => setOptionsHidden(false && isHighlighted)}
			onMouseLeave={() => setOptionsHidden(true && isHighlighted)}
			style={isHighlighted ? { width: '200%' } : {}}
			onClick={(e) => clicked(e)}
		></span>
	);
}

export default ItemCoverIndicator;
