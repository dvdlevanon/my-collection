import { ToggleButton, ToggleButtonGroup } from '@mui/material';
import React from 'react';
import { useQuery } from 'react-query';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';

function TagImageTypeSelector({ tit, onTitChanged }) {
	const titsQuery = useQuery(ReactQueryUtil.TAG_IMAGE_TYPES_KEY, Client.getTagImageTypes);

	return (
		<>
			{titsQuery.isSuccess && (
				<ToggleButtonGroup
					size="small"
					exclusive
					value={tit}
					onChange={(e, newValue) => {
						if (newValue != null) {
							onTitChanged(newValue);
						}
					}}
				>
					{titsQuery.data.map((tit) => {
						return (
							<ToggleButton key={tit.id} value={tit}>
								{tit.nickname}
							</ToggleButton>
						);
					})}
				</ToggleButtonGroup>
			)}
		</>
	);
}

export default TagImageTypeSelector;
