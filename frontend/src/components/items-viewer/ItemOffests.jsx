import { Typography } from '@mui/material';
import React from 'react';
import TimeUtil from '../../utils/time-utils';

function ItemOffests({ item }) {
	const formatStart = () => {
		return TimeUtil.formatDuration(item.start_position ? item.start_position : 0);
	};

	const formatEnd = () => {
		return TimeUtil.formatDuration(item.end_position);
	};

	return (
		<Typography color="bright.darker" variant="caption">
			{formatStart()} - {formatEnd()}
		</Typography>
	);
}

export default ItemOffests;
