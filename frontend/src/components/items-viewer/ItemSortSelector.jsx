import { useTheme } from '@emotion/react';
import { mdiShuffle, mdiSortAscending, mdiSortDescending } from '@mdi/js';
import Icon from '@mdi/react';
import DurationIcon from '@mui/icons-material/Timelapse';
import { Box, Stack, ToggleButton, ToggleButtonGroup, Tooltip } from '@mui/material';
import React from 'react';

function ItemSortSelector({ sortBy, onSortChanged }) {
	const theme = useTheme();

	return (
		<Stack justifyContent="center">
			<Box>
				<ToggleButtonGroup
					size="small"
					exclusive
					value={sortBy}
					onChange={(e, newValue) => {
						if (newValue != null) {
							onSortChanged(newValue);
						}
					}}
					onClick={(e) => {
						onSortChanged(sortBy);
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
					<Tooltip title="Sort by duration" value="duration">
						<ToggleButton value="duration">
							<DurationIcon sx={{ fontSize: theme.iconSize(1) }} />
						</ToggleButton>
					</Tooltip>
				</ToggleButtonGroup>
			</Box>
		</Stack>
	);
}

export default ItemSortSelector;
