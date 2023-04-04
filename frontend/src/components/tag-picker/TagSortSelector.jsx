import { mdiShuffle, mdiSortAscending, mdiSortDescending } from '@mdi/js';
import Icon from '@mdi/react';
import { ToggleButton, ToggleButtonGroup } from '@mui/material';
import React from 'react';

function TagSortSelector({ sortBy, onSortChanged }) {
	return (
		<>
			<ToggleButtonGroup
				size="small"
				exclusive
				value={sortBy}
				onChange={(e, newValue) => {
					if (newValue != null) {
						onSortChanged(newValue);
					}
				}}
			>
				<ToggleButton value="random">
					<Icon path={mdiShuffle} size={1} />
				</ToggleButton>
				<ToggleButton value="title-asc">
					<Icon path={mdiSortAscending} size={1} />
				</ToggleButton>
				<ToggleButton value="title-desc">
					<Icon path={mdiSortDescending} size={1} />
				</ToggleButton>
			</ToggleButtonGroup>
		</>
	);
}

export default TagSortSelector;
