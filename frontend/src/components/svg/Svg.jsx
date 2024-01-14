import { Box } from '@mui/material';
import React, { useEffect, useState } from 'react';
import Client from '../../utils/client';

function Svg({ path, width, height, color }) {
	const [icon, setIcon] = useState('');
	const [viewBox, setViewBox] = useState('');

	useEffect(() => {
		Client.fetchTextFile(path).then((svgText) => {
			if (svgText.startsWith('<svg')) {
				parseSvg(svgText);
			} else {
				console.log('Not an svg ' + path);
			}
		});
	}, [path]);

	const parseSvg = (svgText) => {
		let parser = new DOMParser();
		let svgDoc = parser.parseFromString(svgText, 'image/svg+xml');

		var svgElement = svgDoc.querySelector('svg');
		let viewBox = svgElement.getAttribute('viewBox');
		setViewBox(viewBox);

		let gElement = svgDoc.querySelector('g');
		let gHtml = new XMLSerializer().serializeToString(gElement);
		setIcon(gHtml);
	};

	const getSvg = () => {
		return `<svg viewBox="${viewBox}" width="${width}" height="${height}" fill="${color}">${icon}</svg>`;
	};

	if (!icon) {
		return null;
	}

	return <Box width={width} height={height} dangerouslySetInnerHTML={{ __html: getSvg() }} />;
}

export default Svg;
