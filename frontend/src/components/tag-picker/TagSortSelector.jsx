import { mdiShuffle, mdiSortAscending, mdiSortDescending } from '@mdi/js';
import Icon from '@mdi/react';
import ItemsCountIcon from '@mui/icons-material/Category';
import { ToggleButton, ToggleButtonGroup, Tooltip } from '@mui/material';
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
				<Tooltip title="Do not sort">
					<ToggleButton value="random">
						<Icon path={mdiShuffle} size={1} />
					</ToggleButton>
				</Tooltip>
				<Tooltip title="Sort by title asc">
					<ToggleButton value="title-asc">
						<Icon path={mdiSortAscending} size={1} />
					</ToggleButton>
				</Tooltip>
				<Tooltip title="Sort by title desc">
					<ToggleButton value="title-desc">
						<Icon path={mdiSortDescending} size={1} />
					</ToggleButton>
				</Tooltip>
				<Tooltip title="Sort by items count">
					<ToggleButton value="items-count">
						<ItemsCountIcon />
					</ToggleButton>
				</Tooltip>
			</ToggleButtonGroup>
		</>
	);
}

export default TagSortSelector;
