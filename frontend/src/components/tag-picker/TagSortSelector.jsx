import { useTheme } from '@emotion/react';
import { mdiShuffle, mdiSortAscending, mdiSortDescending } from '@mdi/js';
import Icon from '@mdi/react';
import ItemsCountIcon from '@mui/icons-material/Category';
import { ToggleButton, ToggleButtonGroup, Tooltip } from '@mui/material';
import React from 'react';

function TagSortSelector({ sortBy, onSortChanged }) {
	const theme = useTheme();

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
				<Tooltip title="Do not sort" value="random">
					<ToggleButton value="random">
						<Icon path={mdiShuffle} size={theme.iconSize(1)} />
					</ToggleButton>
				</Tooltip>
				<Tooltip title="Sort by title asc" value="title-asc">
					<ToggleButton value="title-asc">
						<Icon path={mdiSortAscending} size={theme.iconSize(1)} />
					</ToggleButton>
				</Tooltip>
				<Tooltip title="Sort by title desc" value="title-desc">
					<ToggleButton value="title-desc">
						<Icon path={mdiSortDescending} size={theme.iconSize(1)} />
					</ToggleButton>
				</Tooltip>
				<Tooltip title="Sort by items count" value="items-count">
					<ToggleButton value="items-count">
						<ItemsCountIcon sx={{ fontSize: theme.iconSize(1) }} />
					</ToggleButton>
				</Tooltip>
			</ToggleButtonGroup>
		</>
	);
}

export default TagSortSelector;
