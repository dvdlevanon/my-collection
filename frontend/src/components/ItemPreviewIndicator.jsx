import React, { useState } from 'react';
import styles from './ItemPreviewIndicator.module.css';

function ItemPreviewIndicator({ item, preview, isHighlighted }) {
	let [optionsHidden, setOptionsHidden] = useState(true);

	const clicked = (e) => {
		e.stopPropagation();
		console.log('cliecked');
	};

	return (
		<span
			className={styles.preview_indicator + ' ' + (isHighlighted ? styles.selected : styles.unselected)}
			key={preview.id}
			onMouseEnter={() => setOptionsHidden(false && isHighlighted)}
			onMouseLeave={() => setOptionsHidden(true && isHighlighted)}
			style={isHighlighted ? { width: '200%' } : {}}
			onClick={(e) => clicked(e)}
		></span>
	);
}

export default ItemPreviewIndicator;
